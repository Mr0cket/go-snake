package models

import (
	"snake/utils"
)

type Snake struct {
	body []utils.Position
}

func (Snake) New(windowSize utils.WindowDimensions) Snake {
	head := utils.Position{X: windowSize.Cols / 2, Y: windowSize.Rows / 2}

	// Add the new head to the snake
	return Snake{body:[]utils.Position{ head, {X: head.X - 1, Y: head.Y}, {X: head.X - 2, Y: head.Y}}}
}

func (s Snake) Add(utils.Position) {
	s.body[len(s.body)] = s.body[len(s.body)-1]
}

func (s Snake) Move(direction string) {
	// Update the snake's position
	snakeHead := s.body[0]

	switch direction {
	case "up":
		snakeHead.Y--
	case "down":
		snakeHead.Y++
	case "left":
		snakeHead.X--
	case "right":
		snakeHead.X++
	}
	for i := len(s.body) - 1; i > 0; i-- {
	s.body[i] = s.body[i-1]
	}
	s.body[0] = snakeHead
}

func (s Snake) Render() {
	// Render the snake's head
	utils.SetPosition("X", s.body[0].X, s.body[0].Y)
	// Render the snake body
	for _, segment := range s.body[1:] {
		utils.SetPosition("█", segment.X, segment.Y)
	}
}

func (s Snake) HitPoint(point utils.Position) bool {
	head := s.body[0]
	return head.X == point.X && head.Y == point.Y
}

func (s Snake) HitWall(windowSize utils.WindowDimensions) bool {
	head := s.body[0]
	// Check if the snake hit the wall
	return head.X < 1 || head.X > windowSize.Cols-1 || head.Y < 1 || head.Y > windowSize.Rows-1
}

func (s Snake) HitSelf() bool {
	// Check if the snake hit itself
	for _, segment := range s.body[1:] {
		if s.HitPoint(segment) {
			return true
		}
	}
	return false
}
