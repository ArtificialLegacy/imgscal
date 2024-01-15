package statemachine

type CliState int8

const (
	NONE            CliState = -1
	ESRGAN_VERIFY   CliState = 0
	ESRGAN_DOWNLOAD CliState = 1
	ESRGAN_FAIL     CliState = 2
	LANDING_MENU    CliState = 3
	ESRGAN_MANAGE   CliState = 4
	WORKFLOW_MENU   CliState = 5
	WORKFLOW_FINISH CliState = 6
	ESRGAN_X4       CliState = 7
	ESRGAN_ANIMEX4  CliState = 8
)
