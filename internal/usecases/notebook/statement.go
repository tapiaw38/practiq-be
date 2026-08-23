package notebook

import (
	"context"
	"log"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/assistant"
	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

func pageHasImageStatement(contentData string) bool {
	value := strings.TrimSpace(contentData)
	return value != "" && (isLikelyImageData(value) || isImageURL(value))
}

func ensurePageStatement(ctx context.Context, app *appcontext.Context, cfg assistant.Config, page *domain.NotebookPage) {
	if page == nil || strings.TrimSpace(page.StatementText) != "" {
		return
	}
	if !pageHasImageStatement(page.ContentData) {
		return
	}
	if app.Integrations.AssistantGateway == nil || !app.Integrations.AssistantGateway.IsConfigured(cfg) {
		return
	}

	resolved, err := resolveImageForOCR(ctx, app, page.ContentData)
	if err != nil {
		log.Printf("[notebook] statement backfill resolve failed page_id=%s err=%v", page.ID, err)
		return
	}

	transcription, err := app.Integrations.AssistantGateway.AnalyzeNotebookStatement(
		ctx, cfg, normalizeCanvasDataURI(resolved), buildNotebookPromptContext(page),
	)
	if err != nil {
		log.Printf("[notebook] statement backfill failed page_id=%s err=%v", page.ID, err)
		return
	}

	transcription = strings.TrimSpace(transcription)
	if transcription == "" || strings.EqualFold(transcription, "UNREADABLE") {
		return
	}

	page.StatementText = transcription
	if err := app.Repositories.Notebook.UpdatePageStatement(ctx, page.ID, transcription); err != nil {
		log.Printf("[notebook] statement backfill persist failed page_id=%s err=%v", page.ID, err)
	}
}

func transcribePageStatement(ctx context.Context, app *appcontext.Context, teacherID, contentData string, page domain.NotebookPage) string {
	if !pageHasImageStatement(contentData) {
		return ""
	}
	if app.Integrations.AssistantGateway == nil {
		return ""
	}

	profile, _ := app.Repositories.UserProfile.Get(ctx, teacherID)
	cfg := assistant.Config{}
	if profile != nil {
		cfg.BaseURL = profile.AssistantBaseURL
		cfg.APIKey = profile.AssistantAPIKey
	}
	if !app.Integrations.AssistantGateway.IsConfigured(cfg) {
		return ""
	}

	resolved, err := resolveImageForOCR(ctx, app, contentData)
	if err != nil {
		log.Printf("[notebook] statement resolve failed err=%v", err)
		return ""
	}

	transcription, err := app.Integrations.AssistantGateway.AnalyzeNotebookStatement(
		ctx, cfg, normalizeCanvasDataURI(resolved), buildNotebookPromptContext(&page),
	)
	if err != nil {
		log.Printf("[notebook] statement transcription failed err=%v", err)
		return ""
	}

	transcription = strings.TrimSpace(transcription)
	if transcription == "" || strings.EqualFold(transcription, "UNREADABLE") {
		return ""
	}
	return transcription
}
