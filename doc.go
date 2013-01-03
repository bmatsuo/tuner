/*
Tuner controls iTunes via Applescript from the command line.

Examples

Control local iTunes

	tuner play
	tuner next
	tuner rate 5

Or control iTunes on remote machines (this is flakey)

	tuner -host eppc://10.0.0.34 pause

Commands

The following commands are available in tuner

	help    command usage
	info    current track info
	mute    mute/unmute playback
	next    skip to the next track
	open    open iTunes
	pause   pause the current track
	play    start playing the current track
	prev    skip to the previous track
	quit    quit iTunes
	rate    adjust rating of the current track
	status  player status (playing, paused, stopped)
	stop    stop playback
	vol     adjust playback volume


Help

Command usage is available through the "help" command

	tuner help [COMMAND]
*/
package documentation
