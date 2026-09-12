# Suscripciones: cómo funcionan

Qué limita un plan, quién decide el límite y dónde vive cada regla. Escrito
para que alguien que no leyó el código pueda cambiarlo sin romper lo que ya
está resuelto.

## Las tres piezas

| Pieza | Qué sabe | Qué NO sabe |
|---|---|---|
| `practiq-be` (Go) | Qué significa un plan: cuántos alumnos permite, quién cuenta como alumno, cuándo se rechaza uno | Nada de tarjetas ni de cobros |
| `payments-dev` (FastAPI) | Qué se cobra, a quién, cuándo. Guarda planes y suscripciones | **Qué permite un plan.** Lo guarda como metadata opaca y lo devuelve sin leerlo |
| Mercado Pago | La pasarela. Cobra de verdad | Nada de Practiq |

La separación es deliberada: **el significado de un plan lo decide Practiq, no
el servicio de pagos.** Por eso el límite viaja como metadata. Cambiar qué
cuenta como alumno no requiere tocar ni desplegar `payments-dev`, y esa regla
ya cambió una vez (de "por docente" a "por escuela").

## De dónde sale el límite

```
Plan en payments-dev
  metadata: {"max_students": 5, "name": "Inicial"}
        ↓
GET /api/v1/subscriptions/entitlement/{user_id}
  devuelve subscription.plan.plan_metadata tal cual
        ↓
domain.PlanFromMetadata()            internal/domain/subscription.go
  lee "max_students" → TeacherPlan{MaxStudents: 5}
        ↓
TeacherSubscription.CanAddStudent()
  StudentsUsed < MaxStudents
```

Planes reales hoy: Inicial=5, Crecimiento=15, Pro=30.

Si un plan no trae `max_students`, **cae a 1, no a ilimitado**. Un plan mal
cargado no debe vender todo (`PlanFromMetadata`, comentario en el código).

### El plan gratis no existe en Mercado Pago

`domain.FreePlan` está hardcodeado: 1 alumno, 30 días. A propósito — crear una
suscripción que cobra cero para después no cobrarla no compra nada y deja el
primer mes de un docente a merced de que la pasarela esté viva. Pagos entra
cuando entra la plata.

## `scopeFor`: el único lugar donde se resuelve un plan

`internal/usecases/subscription/scope.go`

Firma: `scopeFor(ctx, app, schoolID, teacherID)`.

**El tope sale de la escuela a la que entra el alumno, no del docente.** Todo
docente recibe escuela personal al registrarse (`ensurePersonalSchool` corre en
cada sync de perfil, `sync.go`), así que deducirlo del docente cobraba los
alumnos de una institución contra el plan personal de ese docente — y los
rechazaba cuando ese plan se llenaba, con un tope que la institución no tiene.

`schoolID` vacío significa "la escuela propia del docente". El único caso que
no tiene escuela que nombrar es la asignación directa docente-alumno, que es
justamente el caso personal.

| Estado | Cuándo | Enforcement | Qué muestra la pantalla |
|---|---|---|---|
| `capEnforced` | Escuela `personal` + `billing = subscription` | Sí | El plan y el uso |
| `capNone` | Cualquier `institution`, o cualquier escuela `direct`, o escuela inexistente | No | "Sin límite" (`uncapped: true`) |
| `capUnknown` | `payments-dev` caído | No | El plan gratis |

Una institución no tiene tope **sin importar su `billing`**. Se factura por
contrato; que la columna diga `subscription` no crea un precio por alumno que
nadie acordó. Tampoco hay `max_teachers` ni `max_admins`: sus admins deciden
cuánta gente tiene y nada de eso se vende acá.

El entitlement se lee del **dueño de la escuela** (`school.CreatedBy`), no de
quien está agregando al alumno. Si un admin suma a alguien a la escuela
personal de otro docente, paga el plan del dueño.

**`capNone` y `capUnknown` no son lo mismo aunque los dos dejen pasar a
cualquiera.** Uno es "no hay límite", el otro es "no pude leer el límite".
Colapsarlos hace que en una caída de pagos la pantalla le diga a un docente que
paga que no tiene tope.

Lo consumen tres lugares, y **deben seguir siendo los mismos tres**:

- `allowance.go` → `EnsureCanAddStudent`, el chequeo que rechaza
- `downgrade.go` → preview/apply/reactivate
- `get_mine.go` → la pantalla de suscripción

> Esto era un bug. La pantalla contaba por docente (assignments ∪ enrollments ∪
> grade_memberships) y el enforcement por escuela (`school_members` activos).
> Los dos números discrepaban cada vez que un downgrade desactivaba a alguien o
> una membresía de escuela no se escribía. Si agregás un cuarto consumidor,
> que lea `scopeFor` — no vuelvas a escribir la regla.

## Quién cuenta como alumno

Una sola consulta: `School.CountStudents` →
`school_members WHERE school_id = ? AND role = 'student' AND active`.

### Qué escuela pasa cada camino de alta

Los tres lugares donde un docente gana un alumno. **La escuela que se le pasa al
límite tiene que ser la misma en la que el alumno termina**, o el tope bloquea
usando un contador que ese mismo flujo nunca incrementa:

| Camino | Escuela | Archivo |
|---|---|---|
| Canje de invitación | `invitation.SchoolID` | `student_invitation/redeem.go` |
| Inscripción a curso | `course.SchoolID` | `enrollment/enroll.go` |
| Asignación directa | `""` → la personal del docente | `teacher_student_assignment/assign.go` |

Lo mismo vale para `school.JoinSchool`, que es quien efectivamente crea la
membresía. Si agregás un cuarto camino, pasá la escuela a los dos.

## El mes gratis expira

`domain.EffectiveFreePlan(profileCreatedAt, now)` devuelve **0 alumnos** una vez
pasados los 30 días. No devuelve un flag "expirado".

Ese detalle importa: con cero, los tres consumidores hacen lo correcto sin
saber que existe un trial. El chequeo rechaza al siguiente alumno, el preview
de downgrade lista a los que sobran, la pantalla muestra que no queda nada. Si
en cambio agregás un booleano `expired`, hay que acordarse de mirarlo en cada
lugar — y alguno se va a olvidar.

El mes se cuenta desde `user_profiles.created_at`, que es el único momento que
existe para todo docente.

Pagar termina el efecto del trial: un `entitlement` activo ni consulta la fecha.

## Flujos

### Suscribirse — `POST /api/teachers/me/subscription`

```
Browser: SDK de Mercado Pago tokeniza la tarjeta (public key)
   ↓ card_token_id (un solo uso, expira rápido)
practiq-be: subscribe.go
   - el email del pagador sale de auth-api, NUNCA del body
     (aceptarlo del body deja facturarle a otro)
   - reenvía plan_id + token a payments-dev
   ↓
payments-dev: create_subscription
   - si ya hay una viva y es OTRO plan → la cancela ANTES de crear la nueva
   - crea preapproval en Mercado Pago
```

**La tarjeta nunca toca nuestros servidores.** Solo el token. Eso es lo que
mantiene a Practiq fuera del alcance de PCI.

**El orden cancelar→crear es deliberado.** Al revés, si el segundo paso falla
quedan dos acuerdos vivos y el pagador se lleva dos cobros. Así, lo peor que
pasa es que queda sin plan y puede reintentar. Perder acceso se recupera; un
cobro doble es un reembolso y un cliente menos.

### Pausar / reanudar / cancelar — `POST .../subscription/{pause,resume,cancel}`

`manage_mine.go`. No recibe id de suscripción: la resuelve del docente que
pregunta. No hay id que un llamador pueda cambiar por el de otro.

Busca entre **todas** las suscripciones, no solo la vigente: una pausada no es
un entitlement, y mirando solo entitlements, reanudar sería imposible justo
para quien lo necesita. Las canceladas se saltean (son terminales en la
pasarela).

En la UI, cancelar ofrece pausar primero.

### Bajar de plan — `GET/POST .../subscription/downgrade`

No se aplica solo. El docente ya pagó el período en curso, y los alumnos que
perderían acceso no tomaron la decisión.

- `Preview` dice quiénes quedarían fuera
- `Apply` los desactiva (`school_members.active = false`) — el docente puede
  elegir a cuáles conservar
- `Reactivate` trae a uno de vuelta, con el mismo chequeo que agregar

El orden lo decide `domain.StudentsToDeactivate`: **menos activo primero**
(`student_topic_progress.last_practiced_at`), los que nunca practicaron
primero de todo. Por antigüedad se estaría echando a los que están en clase hoy
y conservando a los que terminaron hace meses.

Desactivar no borra: conservan cuenta e historial.

## Decisiones que parecen bugs y no lo son

**Si `payments-dev` no responde, se deja pasar al alumno.** Tratar "no pude
leer" como "plan gratis" frenaría a todo docente que paga cada vez que el
servicio hipa. Un alumno de más cuesta poco; dejar a un docente sin su clase
cuesta el producto.

**Un alumno que el docente ya tiene siempre pasa.** Los caminos que llaman a
`EnsureCanAddStudent` son idempotentes; reejecutar uno no puede empezar a
fallar porque el plan se llenó mientras tanto. No se está agregando a nadie.

**Editar un plan no recotiza a los ya suscritos.** Cada uno mantiene el acuerdo
al precio que aceptó. Un aumento que se autoaplica es un contracargo, no una
feature.

## Endpoints

| Método | Ruta | Quién |
|---|---|---|
| GET | `/api/public/subscription-plans` | Público (landing). Solo planes activos |
| GET | `/api/teachers/me/subscription` | Docente. Plan + uso + estado |
| GET | `/api/teachers/me/subscription/checkout-config` | Docente. Public key para tokenizar |
| POST | `/api/teachers/me/subscription` | Docente. Suscribirse |
| POST | `/api/teachers/me/subscription/{pause,resume,cancel}` | Docente |
| GET/POST | `/api/teachers/me/subscription/downgrade` | Docente |
| POST | `/api/teachers/me/students/{id}/reactivate` | Docente |
| CRUD | `/api/subscription-plans` | Superadmin |

`payments-dev` **no publica puertos**. Su API autentica a un producto, no a una
persona, así que exponerla dejaría todas las suscripciones al alcance de quien
tenga una clave. `practiq-be` la alcanza por la red interna. La única excepción
es `payments.practiq.com.ar/api/v1/webhooks/*`, porque Mercado Pago no tiene
clave que mandar — firma sus requests y el servicio verifica esa firma.

## Campos de `GET /teachers/me/subscription`

```jsonc
{
  "plan": { "plan_id": 1, "name": "Inicial", "max_students": 5 },
  "active": true,          // paga ahora (no: "le queda mes gratis")
  "status": "authorized",  // palabra de la pasarela: authorized | paused
  "students_used": 3,
  "can_add_student": true,
  "renews_at": "2026-10-11T...",  // fin del período pago, o del mes gratis
  "uncapped": false,       // sin tope: institución o facturación directa
  "trial_expired": false   // el mes gratis se terminó
}
```

## Invariantes — romper esto rompe el producto

1. **Un solo conteo de alumnos.** Todo pasa por `scopeFor` + `School.CountStudents`.
2. **El tope sale de la escuela destino**, nunca del docente. Y esa escuela es la
   misma que recibe la membresía.
3. **Una institución nunca tiene tope**, cualquiera sea su `billing`.
4. **Un plan sin `max_students` permite 1, no infinitos.**
5. **Pagos caído nunca frena a un docente.**
6. **La tarjeta no llega al backend.** Solo el token de un uso.
7. **Cancelar antes de crear** al cambiar de plan.
8. **El email de facturación sale de auth-api**, jamás del request.
9. **Un trial vencido permite 0**, no 1.

## Qué falta (a hoy)

- `MP_WEBHOOK_SECRET` sin configurar: los webhooks se rechazan por firma. Hasta
  que esté, una suscripción no pasa sola de `pending` a `authorized`.
- El access token de Mercado Pago es **productivo**: una prueba cobra de verdad.
- `cancel_due_subscriptions` existe en `payments-dev` pero **nadie lo llama**.
  Requiere un worker programado.

## Archivos

```
practiq-be/
  internal/domain/subscription.go           FreePlan, EffectiveFreePlan, PlanFromMetadata, CanAddStudent
  internal/domain/downgrade.go              StudentsToDeactivate (función pura)
  internal/usecases/subscription/
    scope.go        scopeFor, studentsUsed   ← empezar acá
    allowance.go    EnsureCanAddStudent
    get_mine.go     la pantalla
    downgrade.go    preview/apply/reactivate
    subscribe.go    alta con token de tarjeta
    manage_mine.go  pause/resume/cancel
    plans.go        catálogo + CRUD superadmin
  internal/adapters/web/integrations/payments/payments.go   cliente HTTP

payments-dev/
  src/services/subscription_service.py      toda la lógica de cobro
  src/api/v1/routers/subscriptions.py       endpoints
  src/api/v1/routers/webhooks.py            ingesta firmada de Mercado Pago
```
