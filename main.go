package main

import (
	"context"
	"fmt"
	"os"

	"SungClip/internal/config"
	"SungClip/internal/controllers"
	"SungClip/internal/services"
	"SungClip/internal/types"
	"SungClip/internal/utils"

	"github.com/spf13/cobra"
)

// ── Package-level variables untuk flags ─────────────────────────────
var (
	// Ingest flags
	videoURL string

	// Edit flags
	editTitle         string
	editWidth         int
	editHeight        int
	editCompositionID string
)

// ── Root Command ────────────────────────────────────────────────────
var rootCmd = &cobra.Command{
	Use:   "sungclip",
	Short: "SungClip - AI-powered short video clip generator",
	Long: `SungClip is a CLI tool for creating short video clips 
with word-by-word subtitles from long-form videos.`,
}

var videoIngestionCmd = &cobra.Command{
	Use:   "ingest",
	Short: "Ingest a video from URL for processing",
	Long:  "Download and process a video from a given URL to prepare it for clip generation.",
	RunE:  runVideoIngestion,
}

var videoEditingCmd = &cobra.Command{
	Use:   "edit",
	Short: "Generate edited video clips from ingested content",
	Long:  "Create short video clips with subtitles using the specified composition settings.",
	RunE:  runVideoEditing,
}

// ── Init: register flags & subcommands ──────────────────────────────
func init() {
	// Ingest flags
	videoIngestionCmd.Flags().StringVarP(&videoURL, "url", "u", "", "Video URL to ingest (required)")
	_ = videoIngestionCmd.MarkFlagRequired("url")

	// Edit flags
	videoEditingCmd.Flags().StringVarP(&editTitle, "title", "t", "", "Title for the generated clips (required)")
	videoEditingCmd.Flags().IntVarP(&editWidth, "width", "W", 1080, "Video width in pixels")
	videoEditingCmd.Flags().IntVarP(&editHeight, "height", "H", 1920, "Video height in pixels")
	videoEditingCmd.Flags().StringVarP(&editCompositionID, "comp", "c", "", "Composition ID from ingestion (required)")
	_ = videoEditingCmd.MarkFlagRequired("title")
	_ = videoEditingCmd.MarkFlagRequired("comp")

	// Add subcommands
	rootCmd.AddCommand(videoIngestionCmd)
	rootCmd.AddCommand(videoEditingCmd)
}

// ── Main ────────────────────────────────────────────────────────────
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// ── Handler Functions ───────────────────────────────────────────────
func runVideoIngestion(cmd *cobra.Command, args []string) error {
	// Setup dependencies
	cfg := config.NewConfig()
	u := utils.NewUtils(cfg)
	svc := services.NewServices(u)
	ctrl := controllers.NewControllers(u, svc)

	ctx := context.Background()

	req := &types.RequestVideoIngestion{
		VideoURL: videoURL,
	}

	resp, err := ctrl.VideoIngestion(ctx, req)
	if err != nil {
		return fmt.Errorf("ingestion failed: %w", err)
	}

	fmt.Printf("✅ Ingestion successful!\n")
	fmt.Printf("   Title:      %s\n", resp.Title)
	fmt.Printf("   PromptPath: %s\n", resp.PromptPath)

	return nil
}

func runVideoEditing(cmd *cobra.Command, args []string) error {
	// Setup dependencies
	cfg := config.NewConfig()
	u := utils.NewUtils(cfg)
	svc := services.NewServices(u)
	ctrl := controllers.NewControllers(u, svc)

	ctx := context.Background()

	req := &types.RequestVideoEditing{
		Title:         editTitle,
		Width:         editWidth,
		Height:        editHeight,
		CompositionID: editCompositionID,
	}

	resp, err := ctrl.VideoEditing(ctx, req)
	if err != nil {
		return fmt.Errorf("editing failed: %w", err)
	}

	fmt.Printf("✅ Editing successful!\n")
	fmt.Printf("   Title:           %s\n", resp.Title)
	fmt.Printf("   Total Clips:     %d\n", resp.TotalClips)
	fmt.Printf("   Result Path:     %s\n", resp.ResultVideoPath)

	return nil
}