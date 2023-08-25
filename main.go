package main

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
)

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range text {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func notify(s tcell.Screen, style tcell.Style, text string) {
	width, height := s.Size()
	text = fmt.Sprintf(fmt.Sprintf("%%-%ds", width), text)
	drawText(s, 0, height - 1, width - 1, height - 1, style, text)
}

func main() {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	// boxStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorPurple)

	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	// Here's how to get the screen size when you need it.
	// xmax, ymax := s.Size()

	// Here's an example of how to inject a keystroke where it will
	// be picked up by the next PollEvent call.  Note that the
	// queue is LIFO, it has a limited length, and PostEvent() can
	// return an error.
	// s.PostEvent(tcell.NewEventKey(tcell.KeyRune, rune('a'), 0))

	// Event loop
	for {
		// Update screen
		s.Show()

		// Poll event
		ev := s.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventError:
			notify(s, defStyle, fmt.Sprintf("Error: %+v", ev))
		case *tcell.EventInterrupt:
			notify(s, defStyle, fmt.Sprintf("Interrupt: %+v", ev))
		case *tcell.EventPaste:
			notify(s, defStyle, fmt.Sprintf("Paste: %+v", ev))
		case *tcell.EventTime:
			notify(s, defStyle, fmt.Sprintf("Time: %+v", ev))
		case *tcell.EventMouse:
			notify(s, defStyle, fmt.Sprintf("Mouse: %+v", ev))
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyCtrlC {
				return
			}
			notify(s, defStyle, fmt.Sprintf("Key: %s [%+v]", ev.Name(), ev))
		default:
			notify(s, defStyle, fmt.Sprintf("Other? %+v", ev))
		}
	}
}
