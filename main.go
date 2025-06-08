package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	screenWidth  = 640
	screenHeight = 480
	fontSize     = 24
	kanjiSize    = 72
)

// SushiItem represents a sushi item with its kanji and English name
type SushiItem struct {
	Kanji  string
	English string
}

// Game represents the main game state
type Game struct {
	sushiItems     []SushiItem
	currentItem    SushiItem
	options        []string
	correctIndex   int
	score          int
	bestScore      int
	gameOver       bool
	mplusNormalFont font.Face
	mplusBigFont    font.Face
	smallerFont     font.Face
	rand           *rand.Rand
	tt             *opentype.Font
}

// Initialize the game
func NewGame() *Game {
	g := &Game{
		sushiItems: []SushiItem{
			{"鮪", "Tuna"},
			{"鮭", "Salmon"},
			{"鰤", "Yellowtail"},
			{"鯛", "Sea Bream"},
			{"鰹", "Bonito"},
			{"鱈", "Cod"},
			{"鰻", "Eel"},
			{"鱚", "Japanese Whiting"},
			{"鯖", "Mackerel"},
			{"鯵", "Horse Mackerel"},
			{"鮟鱇", "Anglerfish"},
			{"鱧", "Pike Conger"},
			{"鱸", "Sea Bass"},
			{"鰆", "Spanish Mackerel"},
			{"鰈", "Flounder"},
			{"鰺", "Amberjack"},
			{"鱒", "Trout"},
			{"鰯", "Sardine"},
		},
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	// Load fonts
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	
	g.tt = tt

	g.mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	g.mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    kanjiSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	
	g.smallerFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSize * 0.8,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	g.startNewRound()
	return g
}

// Start a new round with a random sushi item and options
func (g *Game) startNewRound() {
	// Select a random sushi item
	itemIndex := g.rand.Intn(len(g.sushiItems))
	g.currentItem = g.sushiItems[itemIndex]

	// Create options (one correct, two incorrect)
	g.options = make([]string, 3)
	g.correctIndex = g.rand.Intn(3)
	g.options[g.correctIndex] = g.currentItem.English

	// Fill other options with random incorrect answers
	usedIndices := map[int]bool{itemIndex: true}
	for i := 0; i < 3; i++ {
		if i == g.correctIndex {
			continue
		}

		// Find a random item that's not already used
		var randomIndex int
		for {
			randomIndex = g.rand.Intn(len(g.sushiItems))
			if !usedIndices[randomIndex] {
				break
			}
		}
		usedIndices[randomIndex] = true
		g.options[i] = g.sushiItems[randomIndex].English
	}
	
	// Debug output to verify the correct answer
	fmt.Printf("New round: Kanji=%s, Correct=%s (index %d)\n", 
		g.currentItem.Kanji, g.currentItem.English, g.correctIndex)
}

// Update the game state
func (g *Game) Update() error {
	if g.gameOver {
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.gameOver = false
			g.score = 0
			g.startNewRound()
		}
		return nil
	}

	// Check for mouse clicks on option cards
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		
		// Check if click is on any of the option cards
		cardWidth := 180
		cardHeight := 80
		cardSpacing := 20
		totalCardsWidth := 3*cardWidth + 2*cardSpacing
		startX := (screenWidth - totalCardsWidth) / 2
		
		for i := 0; i < 3; i++ {
			cardX := startX + i*(cardWidth+cardSpacing)
			cardY := screenHeight - 150
			
			if x >= cardX && x < cardX+cardWidth && y >= cardY && y < cardY+cardHeight {
				// Option selected
				if i == g.correctIndex {
					// Correct answer
					g.score++
					if g.score > g.bestScore {
						g.bestScore = g.score
					}
					g.startNewRound()
				} else {
					// Wrong answer
					g.gameOver = true
				}
				return nil // Return after processing the click
			}
		}
	}

	return nil
}

// Draw the game
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw background (light wood color for sushi restaurant feel)
	screen.Fill(color.RGBA{240, 220, 180, 255})
	
	// Draw decorative elements to simulate sushi restaurant
	ebitenutil.DrawRect(screen, 0, 0, screenWidth, 40, color.RGBA{80, 40, 20, 255})
	ebitenutil.DrawRect(screen, 0, screenHeight-40, screenWidth, 40, color.RGBA{80, 40, 20, 255})
	
	if g.gameOver {
		// Game over screen
		msg := "Game Over!"
		bounds := text.BoundString(g.mplusNormalFont, msg)
		x := (screenWidth - bounds.Dx()) / 2
		y := screenHeight/2 - 40
		text.Draw(screen, msg, g.mplusNormalFont, x, y, color.RGBA{200, 30, 30, 255})
		
		scoreMsg := fmt.Sprintf("Your Score: %d", g.score)
		bounds = text.BoundString(g.mplusNormalFont, scoreMsg)
		x = (screenWidth - bounds.Dx()) / 2
		y += 40
		text.Draw(screen, scoreMsg, g.mplusNormalFont, x, y, color.RGBA{0, 0, 0, 255})
		
		bestMsg := fmt.Sprintf("Best Score: %d", g.bestScore)
		bounds = text.BoundString(g.mplusNormalFont, bestMsg)
		x = (screenWidth - bounds.Dx()) / 2
		y += 40
		text.Draw(screen, bestMsg, g.mplusNormalFont, x, y, color.RGBA{0, 0, 0, 255})
		
		restartMsg := "Press SPACE to restart"
		bounds = text.BoundString(g.mplusNormalFont, restartMsg)
		x = (screenWidth - bounds.Dx()) / 2
		y += 60
		text.Draw(screen, restartMsg, g.mplusNormalFont, x, y, color.RGBA{0, 0, 0, 255})
		
		return
	}
	
	// Draw score
	scoreText := fmt.Sprintf("Score: %d  Best: %d", g.score, g.bestScore)
	text.Draw(screen, scoreText, g.mplusNormalFont, 20, 30, color.RGBA{255, 255, 255, 255})
	
	// Draw the kanji character
	kanjiText := g.currentItem.Kanji
	bounds := text.BoundString(g.mplusBigFont, kanjiText)
	x := (screenWidth - bounds.Dx()) / 2
	y := screenHeight/3 + bounds.Dy()/2
	text.Draw(screen, kanjiText, g.mplusBigFont, x, y, color.RGBA{0, 0, 0, 255})
	
	// Draw the options as cards
	cardWidth := 180
	cardHeight := 80
	cardSpacing := 20
	totalCardsWidth := 3*cardWidth + 2*cardSpacing
	startX := (screenWidth - totalCardsWidth) / 2
	
	for i := 0; i < 3; i++ {
		cardX := startX + i*(cardWidth+cardSpacing)
		cardY := screenHeight - 150
		
		// Draw card background
		ebitenutil.DrawRect(screen, float64(cardX), float64(cardY), float64(cardWidth), float64(cardHeight), color.RGBA{220, 220, 220, 255})
		ebitenutil.DrawRect(screen, float64(cardX+2), float64(cardY+2), float64(cardWidth-4), float64(cardHeight-4), color.RGBA{255, 255, 255, 255})
		
		// Draw option text - handle long text with wrapping if needed
		optionText := g.options[i]
		optionBounds := text.BoundString(g.mplusNormalFont, optionText)
		
		// If text is too wide for the card, use a smaller font
		if optionBounds.Dx() > cardWidth-20 {
			// Use the smaller font for this text
			optionBounds = text.BoundString(g.smallerFont, optionText)
			optionX := cardX + (cardWidth-optionBounds.Dx())/2
			optionY := cardY + cardHeight/2 + optionBounds.Dy()/4
			text.Draw(screen, optionText, g.smallerFont, optionX, optionY, color.RGBA{0, 0, 0, 255})
		} else {
			// Normal text rendering for text that fits
			optionX := cardX + (cardWidth-optionBounds.Dx())/2
			optionY := cardY + cardHeight/2 + optionBounds.Dy()/4
			text.Draw(screen, optionText, g.mplusNormalFont, optionX, optionY, color.RGBA{0, 0, 0, 255})
		}
	}
	
	// Draw instructions
	instructionText := "Click on the correct English name"
	bounds = text.BoundString(g.mplusNormalFont, instructionText)
	x = (screenWidth - bounds.Dx()) / 2
	y = screenHeight - 60
	text.Draw(screen, instructionText, g.mplusNormalFont, x, y, color.RGBA{80, 40, 20, 255})
}

// Layout implements ebiten.Game's Layout.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Sushi Kanji Quiz")
	
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
