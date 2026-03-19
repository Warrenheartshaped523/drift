Output demo/drift.gif

Set Shell "zsh"
Set FontSize 15
Set Width 1400
Set Height 800
Set Padding 20
Set Theme "Catppuccin Mocha"
# Note: vhs theme names — use "Catppuccin Mocha", "Dracula", "TokyoNight", "nord"
Set TypingSpeed 80ms
Set PlaybackSpeed 1

# Give the shell a moment to settle
Sleep 500ms

# Waveform — the most visually striking scene
Type "drift --scene waveform --theme catppuccin --duration 0"
Enter
Sleep 8s
Enter

Sleep 400ms

# Constellation
Type "drift --scene constellation --theme cosmic --duration 0"
Enter
Sleep 8s
Enter

Sleep 400ms

# Rain
Type "drift --scene rain --theme dracula --duration 0"
Enter
Sleep 6s
Enter

Sleep 400ms

# Particles
Type "drift --scene particles --theme nord --duration 0"
Enter
Sleep 6s
Enter

Sleep 300ms
