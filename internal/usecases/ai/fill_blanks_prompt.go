package ai

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func buildFillBlanksPuzzleContext(exercise *domain.Exercise) string {
	if exercise == nil {
		return ""
	}

	ids, err := domain.BlankIDsInStatement(exercise.Question)
	if err != nil {
		return ""
	}
	statement := exercise.Question
	for _, id := range ids {
		statement = strings.ReplaceAll(statement, fmt.Sprintf("{{%d}}", id), fmt.Sprintf("[hueco %d]", id))
	}

	var metadata domain.FillBlanksMetadata
	_ = json.Unmarshal([]byte(exercise.Metadata), &metadata)

	var sb strings.Builder
	sb.WriteString("- Tipo: rompecabezas de completar huecos. El alumno toca un bloque y luego el hueco.\n")
	sb.WriteString("- Enunciado: ")
	sb.WriteString(statement)
	sb.WriteString("\n- Huecos: ")
	for index, id := range ids {
		if index > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%d", id))
	}
	if len(metadata.Options) > 0 {
		sb.WriteString("\n- Bloques visibles: ")
		sb.WriteString(strings.Join(metadata.Options, ", "))
	}
	sb.WriteString("\n")
	return sb.String()
}

func formatFillBlanksAnswer(raw string) string {
	answers, err := domain.ParseBlanksAnswer(raw)
	if err != nil || len(answers) == 0 {
		return "sin bloques colocados"
	}

	keys := make([]string, 0, len(answers))
	for key := range answers {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("Hueco %s: %s", key, answers[key]))
	}
	return strings.Join(parts, "; ")
}
