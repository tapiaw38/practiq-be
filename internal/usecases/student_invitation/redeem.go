package studentinvitation

import (
	"context"
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/invitecode"
)

type (
	RedeemUsecase interface {
		Execute(ctx context.Context, studentID, rawCode string) (*RedeemOutput, apperrors.ApplicationError)
	}

	redeemUsecase struct {
		contextFactory appcontext.Factory
	}

	RedeemOutput struct {
		Data RedeemData `json:"data"`
	}

	RedeemData struct {
		TeacherID   string `json:"teacher_id"`
		TeacherName string `json:"teacher_name"`
		// AlreadyLinked distingue el canje nuevo del repetido, para que la
		// pantalla no anuncie dos veces lo mismo.
		AlreadyLinked bool `json:"already_linked"`
	}
)

func NewRedeemUsecase(contextFactory appcontext.Factory) RedeemUsecase {
	return &redeemUsecase{contextFactory: contextFactory}
}

func (u *redeemUsecase) Execute(ctx context.Context, studentID, rawCode string) (*RedeemOutput, apperrors.ApplicationError) {
	app := u.contextFactory()
	now := time.Now()

	if !limiter.allow(studentID, now) {
		return nil, apperrors.NewApplicationError(mappings.InvitationTooManyAttemptsError, nil)
	}

	code := invitecode.Normalize(rawCode)
	if len(code) != invitecode.Length {
		limiter.recordFailure(studentID, now)
		return nil, apperrors.NewApplicationError(mappings.InvitationInvalidCodeError, nil)
	}

	// El perfil dice si es alumno; el rol del token no alcanza porque el
	// administrador también entra por acá y no tiene sentido que se asigne a
	// sí mismo como alumno de alguien.
	profile, err := app.Repositories.UserProfile.Get(ctx, studentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	if profile == nil || profile.ProfileType != "student" {
		return nil, apperrors.NewApplicationError(mappings.InvitationTeacherRedeemError, nil)
	}

	invitation, err := app.Repositories.StudentInvitation.GetByCode(ctx, code)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.InvitationGetError, err)
	}
	if invitation == nil || !invitation.IsUsable(now) {
		limiter.recordFailure(studentID, now)
		return nil, apperrors.NewApplicationError(mappings.InvitationInvalidCodeError, nil)
	}

	firstTime, err := app.Repositories.StudentInvitation.Redeem(ctx, invitation.ID, studentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.InvitationRedeemError, err)
	}

	// Assign es idempotente y reactiva el vínculo si estaba inactivo, así que
	// repetir el canje no rompe nada.
	if err := app.Repositories.TeacherStudentAssignment.Assign(ctx, domain.TeacherStudentAssignment{
		TeacherID: invitation.TeacherID,
		StudentID: studentID,
		Status:    "active",
	}); err != nil {
		return nil, apperrors.NewApplicationError(mappings.InvitationRedeemError, err)
	}

	limiter.clear(studentID)

	teacherName := ""
	if teacher, err := app.Repositories.UserProfile.Get(ctx, invitation.TeacherID); err == nil && teacher != nil {
		teacherName = teacher.Name
	}

	return &RedeemOutput{Data: RedeemData{
		TeacherID:     invitation.TeacherID,
		TeacherName:   teacherName,
		AlreadyLinked: !firstTime,
	}}, nil
}
