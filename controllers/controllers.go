package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"study-buddy-backend/models"
	"study-buddy-backend/services/db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	// --- FIX: Corrected and added necessary imports for the Gemini API ---
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type CreateModuleRequest struct {
	LessonID   string            `json:"lesson_id"`
	ModuleType models.ModuleType `json:"module_type"`
	ModuleData json.RawMessage   `json:"module_data"` // raw JSON, like your nested object
	ModuleName string `json:"module_name"`
	ModuleDesc string `json:"module_description"`
}

type ModuleAllRequest struct {
	LessonId string `json:"lesson_id" binding:"required"`
}

func CreateModule(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateModuleRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := db.CreateModule(dbConn, req.LessonID, string(req.ModuleType), req.ModuleData , req.ModuleName, req.ModuleDesc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to create module in database",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Module created successfully",
		})
	}
}

func GetAllModule(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ModuleAllRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		modules, err := db.FindAllModule(dbConn, req.LessonId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"lessons": modules,
		})
	}
}

func GenerateVnovel(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Load .env to get API key
		_ = godotenv.Load()
		apiKey := os.Getenv("GEMAPI")
		if apiKey == "" {
			// FIX: Corrected error message to match the key being checked
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing GEMAPI key in .env"})
			return
		}

		// Parse prompt from request
		var req struct {
			Prompt string `json:"prompt"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Prompt) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing prompt in request body"})
			return
		}

		ctx := context.Background()
		// --- FIX 1: Use `option.WithAPIKey` for client initialization ---
		client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
		if err != nil {
			log.Println("Gemini client init error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Gemini client"})
			return
		}
		// --- FIX 2: `client.Close()` is deprecated and removed ---
		// No defer statement is needed for the new client.

		model := client.GenerativeModel("gemini-2.5-pro")
		fullPrompt := formatPrompt(req.Prompt)

		// --- FIX 3: Use `GenerateContent` and `genai.Text` ---
		resp, err := model.GenerateContent(ctx, genai.Text(fullPrompt))
		if err != nil {
			log.Println("Gemini generation error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate content"})
			return
		}

		// --- FIX 4: Correctly parse the response content ---
		var textContent string
		if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
			if txt, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
				textContent = string(txt)
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected response format from API"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No content generated"})
			return
		}

		// Extract and parse markdown-wrapped JSON from Gemini response
		clean := extractJSON(textContent)

		var jsonData any
		if err := json.Unmarshal([]byte(clean), &jsonData); err != nil {
			log.Println("Failed to parse JSON:", err, ". Raw content:", clean)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid JSON from Gemini", "raw": clean})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":        "success",
			"canvasContent": jsonData,
		})
	}
}

func GetModulesByModuleID(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		moduleID := c.Param("module_id")
		if moduleID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "module_id is required in path"})
			return
		}

		modules, err := db.GetModulesByModuleID(dbConn, moduleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve modules", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"modules": modules,
		})
	}
}

func formatPrompt(userPrompt string) string {
	// Replaced the old prompt with the new, more detailed one.
	// Also corrected the inconsistency between `MapsTo` and `MapsTo`.
	return fmt.Sprintf(`
You are an expert interactive scene designer. Your task is to generate a JSON structure representing a multi-frame interactive story based on the user's prompt.

The output MUST be a JSON array of frames. Each frame is an array of element objects that compose the scene.

---

### **1. Core Data Structure**

Your entire output must be a JSON array of frames (an array of arrays), like this: `+"`Element[][]`"+`

The structure for each element object is as follows. You must strictly adhere to these property names and types.

`+"```typescript"+`
// This is the structure you must follow.
// Do not include this TypeScript code in your JSON output.
// The base properties for all elements:
type BaseElement = {  
  id: string;     // Unique identifier for the element within a frame
  x: number;        // The x-coordinate of the top-left corner
  y: number;        // The y-coordinate of the top-left corner
  width: number;
  height: number;
}

// Specific element types:
type RectElement = BaseElement & { type: 'rect'; }  
type TextElement = BaseElement & { type: 'text'; text: string; fontSize: number; }  
type ButtonElement = BaseElement & { type: 'button'; text: string; navigateTo: number; }  
type ImageElement = BaseElement & { type: 'image'; src: string; }
`+"```"+`

---

### **2. Example of a Perfect Response**

Here is a two-frame example. The first frame has a character, text, and a button. The button uses `+"`MapsTo: 1`"+` to link to the second frame (at index 1).

`+"```json"+`
[
  [
    {
      "id": "char1_frame1",
      "type": "image",
      "src": "teacher_1.svg",
      "x": 400,
      "y": 160,
      "width": 150,
      "height": 200
    },
    {
      "id": "text_frame1",
      "type": "text",
      "text": "Ready to begin the quiz?",
      "fontSize": 18,
      "x": 60,
      "y": 80,
      "width": 280,
      "height": 40
    },
    {
      "id": "start_button",
      "type": "button",
      "text": "Yes!",
      "navigateTo": 1,
      "x": 150,
      "y": 150,
      "width": 100,
      "height": 30
    }
  ],
  [
    {
      "id": "char2_frame2",
      "type": "image",
      "src": "student_1.svg",
      "x": 100,
      "y": 160,
      "width": 150,
      "height": 200
    },
    {
      "id": "text_frame2",
      "type": "text",
      "text": "Great! Here is the first question.",
      "fontSize": 20,
      "x": 100,
      "y": 80,
      "width": 400,
      "height": 50
    }
  ]
]
`+"```"+`

---

### **3. Constraints & Rules**

You MUST adhere strictly to these rules:

1.  **Canvas Boundaries**: The canvas is **800 units wide** and **480 units high**. The coordinate system starts with (0,0) at the top-left corner. All elements' positions (`+"`x`"+`, `+"`y`"+`) and dimensions (`+"`width`"+`, `+"`height`"+`) must ensure the element stays entirely within these bounds.
2.  **Layout**: Avoid overlapping interactive elements like buttons or important text. Ensure elements are logically placed on the screen.
3.  **Image Grounding**: For character images, their `+"`y`"+` coordinate should be high enough (e.g., > 150) so they appear grounded at the bottom of the scene, not floating in the middle.
4.  **Available Images**: You MUST only use the following file names for the `+"`src`"+` property of 'image' elements: `+"`'student_1.svg'`"+`, `+"`'student_2.svg'`"+`, `+"`'teacher_1.svg'`"+`, `+"`'teacher_2.svg'`"+`, `+"`'teacher_3.svg'`"+`. Do not create new image names.
5.  **Brevity**: Keep the `+"`text`"+` property for 'text' and 'button' elements short and concise, preferably under 10 words.
6.  **Navigation**: The `+"`MapsTo`"+` property in 'button' elements MUST be the **0-based integer index** of the frame it links to. For example, to go to the first frame, use `+"`0`"+`; for the second frame, use `+"`1`"+`.
7.  **Unique IDs**: Within a single frame (an inner array), every element's `+"`id`"+` must be a unique string.
8.  **Strict JSON Format**: Your entire response MUST be a single, valid JSON array of arrays, wrapped in a markdown code block that starts with `+"```json`"+` and ends with `+"```"+`. Do not include any explanatory text outside this code block.

---

Now, generate a new interactive story with at least 3 frames based on the following topic: "%s"
`, userPrompt)
}

func extractJSON(text string) string {
	// Remove the ```json ... ``` block
	if strings.HasPrefix(text, "```json") {
		text = strings.TrimPrefix(text, "```json")
	}
	if strings.HasSuffix(text, "```") {
		text = strings.TrimSuffix(text, "```")
	}
	return strings.TrimSpace(text)
}

type EditModuleRequest struct {
	ModuleId   string          `json:"module_id"`
	ModuleData json.RawMessage `json:"module_data"`
	LessonId   string          `json:"lesson_id"`
}

func EditModuleData(dbConn *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req EditModuleRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := db.EditModule(dbConn, req.ModuleId, req.ModuleData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to create module in database",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Module Edited successfully",
		})
	}
}
