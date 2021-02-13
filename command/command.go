package command

import (
	"github.com/singalhimanshu/taskgo/parser"
)

type Command interface {
	Do(*parser.Data) error
	Undo(*parser.Data) error
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
	return c.history[c.history_position].Do(c.data)
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
	if _, err := data.RemoveTask(a.listIdx, a.taskPos); err != nil {
		return err
	}
	return nil
}

type RemoveTaskCommand struct {
	listIdx   int
	taskTitle string
	taskDesc  string
	taskPos   int
}

func CreateRemoveTaskCommand(listIdx, taskPos int) *RemoveTaskCommand {
	return &RemoveTaskCommand{
		listIdx: listIdx,
		taskPos: taskPos,
	}
}

func (r *RemoveTaskCommand) Do(data *parser.Data) error {
	taskData, err := data.RemoveTask(r.listIdx, r.taskPos)
	if err != nil {
		return err
	}
	r.taskTitle = taskData.ItemName
	r.taskDesc = taskData.ItemDescription
	return nil
}

func (r *RemoveTaskCommand) Undo(data *parser.Data) error {
	return data.AddNewTask(r.listIdx, r.taskTitle, r.taskDesc, r.taskPos)
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
