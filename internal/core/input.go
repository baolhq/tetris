package core

import "github.com/hajimehoshi/ebiten/v2"

type KeyState struct {
	IsDown      bool
	WasPressed  bool
	WasReleased bool
}

type Action string

const (
	InputUp    Action = "MoveUp"
	InputDown  Action = "MoveDown"
	InputLeft  Action = "MoveLeft"
	InputRight Action = "MoveRight"
	InputPause Action = "Pause"
	InputEnter Action = "Enter"
)

type InputManager struct {
	keyStates map[ebiten.Key]*KeyState
	keyMap    map[Action][]ebiten.Key
}

var Input = &InputManager{
	keyStates: make(map[ebiten.Key]*KeyState),
	keyMap:    make(map[Action][]ebiten.Key),
}

func init() {
	Input.RegisterAction(InputUp, ebiten.KeyUp, ebiten.KeyW)
	Input.RegisterAction(InputDown, ebiten.KeyDown, ebiten.KeyS)
	Input.RegisterAction(InputLeft, ebiten.KeyLeft, ebiten.KeyA)
	Input.RegisterAction(InputRight, ebiten.KeyRight, ebiten.KeyD)
	Input.RegisterAction(InputPause, ebiten.KeyEscape)
	Input.RegisterAction(InputEnter, ebiten.KeyEnter)
}

func (i *InputManager) RegisterAction(action Action, keys ...ebiten.Key) {
	i.keyMap[action] = keys
	for _, key := range keys {
		if _, exists := i.keyStates[key]; !exists {
			i.keyStates[key] = &KeyState{}
		}
	}
}

func (i *InputManager) Update() {
	for key, state := range i.keyStates {
		pressed := ebiten.IsKeyPressed(key)
		state.WasPressed = !state.IsDown && pressed
		state.WasReleased = state.IsDown && !pressed
		state.IsDown = pressed
	}
}

func (i *InputManager) WasPressed(action Action) bool {
	return i.checkKeys(i.keyMap[action], func(s *KeyState) bool { return s.WasPressed })
}

func (i *InputManager) WasReleased(action Action) bool {
	return i.checkKeys(i.keyMap[action], func(s *KeyState) bool { return s.WasReleased })
}

func (i *InputManager) IsDown(action Action) bool {
	return i.checkKeys(i.keyMap[action], func(s *KeyState) bool { return s.IsDown })
}

func (i *InputManager) checkKeys(keys []ebiten.Key, checker func(*KeyState) bool) bool {
	for _, key := range keys {
		if state, ok := i.keyStates[key]; ok && checker(state) {
			return true
		}
	}
	return false
}
