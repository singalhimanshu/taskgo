package command

import (
	"github.com/singalhimanshu/taskgo/parser"
)

type Command interface {
	Do(*parser.Data) error
	Undo(*parser.Data) error
	Redo(*parser.Data) error
}

type CommandManager struct {
	history          []Command
	history_position int
	data             *parser.Data
}

func CreateNewCommand(data *parser.Data) *CommandManager {
	return &CommandManager{
		history:          []Command{CreateEmptyCommand()},
		history_position: 0,
		data:             data,
	}
}

func (c *CommandManager) Execute(command Command) error {
	if len(c.history) != c.history_position+1 {
		c.history = c.history[:c.history_position+1]
	}
	err := command.Do(c.data)
	if err != nil {
		return err
	}
	c.history = append(c.history, command)
	c.history_position++
	return nil
}

func (c *CommandManager) Undo() error {
	if c.history_position == 0 {
		return nil
	}
	err := c.history[c.history_position].Undo(c.data)
	if err != nil {
		return err
	}
	c.history_position--
	return nil
}

func (c *CommandManager) Redo() error {
	if c.history_position+1 == len(c.history) {
		return nil
	}
	c.history_position++
	return c.history[c.history_position].Redo(c.data)
}

type AddTaskCommand struct {
	listIdx   int
	taskTitle string
	taskDesc  string
	taskPos   int
}

func CreateAddTaskCommand(listIdx int, taskTitle, taskDesc string, taskPos int) *AddTaskCommand {
	return &AddTaskCommand{
		listIdx:   listIdx,
		taskTitle: taskTitle,
		taskDesc:  taskDesc,
		taskPos:   taskPos,
	}
}

func (a *AddTaskCommand) Do(data *parser.Data) error {
	return data.AddNewTask(a.listIdx, a.taskTitle, a.taskDesc, a.taskPos)
}

func (a *AddTaskCommand) Undo(data *parser.Data) error {
	return data.RemoveTask(a.listIdx, a.taskPos)
}

func (a *AddTaskCommand) Redo(data *parser.Data) error {
	return a.Do(data)
}

type SwapListItemCommand struct {
	listIdx       int
	taskIdxFirst  int
	taskIdxSecond int
}

func CreateSwapListItemCommand(listIdx, taskIdxFirst, taskIdxSecond int) *SwapListItemCommand {
	return &SwapListItemCommand{
		listIdx:       listIdx,
		taskIdxFirst:  taskIdxFirst,
		taskIdxSecond: taskIdxSecond,
	}
}

func (s *SwapListItemCommand) Do(data *parser.Data) error {
	return data.SwapListItems(s.listIdx, s.taskIdxFirst, s.taskIdxSecond)
}

func (s *SwapListItemCommand) Undo(data *parser.Data) error {
	return data.SwapListItems(s.listIdx, s.taskIdxSecond, s.taskIdxFirst)
}

func (s *SwapListItemCommand) Redo(data *parser.Data) error {
	return s.Do(data)
}

type MoveTaskCommand struct {
	prevTaskIdx int
	prevListIdx int
	newListIdx  int
}

func CreateMoveTaskCommand(prevTaskIdx, prevListIdx, newListIdx int) *MoveTaskCommand {
	return &MoveTaskCommand{
		prevTaskIdx: prevTaskIdx,
		prevListIdx: prevListIdx,
		newListIdx:  newListIdx,
	}
}

func (s *MoveTaskCommand) Do(data *parser.Data) error {
	return data.MoveTask(s.prevTaskIdx, s.prevListIdx, s.newListIdx)
}

func (s *MoveTaskCommand) Undo(data *parser.Data) error {
	return data.MoveTask(s.prevTaskIdx, s.newListIdx, s.prevListIdx)
}

func (s *MoveTaskCommand) Redo(data *parser.Data) error {
	return s.Do(data)
}

type EmptyCommand struct {
}

func CreateEmptyCommand() *EmptyCommand {
	return &EmptyCommand{}
}

func (e EmptyCommand) Do(data *parser.Data) error {
	return nil
}

func (e EmptyCommand) Undo(data *parser.Data) error {
	return nil
}

func (e EmptyCommand) Redo(data *parser.Data) error {
	return nil
}
