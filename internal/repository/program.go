package repository

type IProgramRepository interface {
}

type ProgramRepository struct {
	// attach db here
}

func (p ProgramRepository) InsertProgram()
