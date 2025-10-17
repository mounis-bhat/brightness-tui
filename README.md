# 🌞 Brightness TUI

A beautiful, interactive terminal user interface (TUI) for controlling screen brightness on Linux systems. Built with Go and the Bubble Tea framework.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)

## ✨ Features

- 🎨 **Beautiful UI**: Elegant, centered interface with rounded borders and smooth animations
- 󰛨 **Dynamic Brightness Icons**: Clean brightness icon that changes based on brightness level (󱩎 to 󰛨)
  - 10 different brightness levels with corresponding icons
  - Consistent size prevents window resizing
- ⌨️ **Multiple Control Methods**: 
  - Arrow keys or vim keys (↑/↓, k/j) for fine-tuned adjustments (±5%)
  - Number keys (1-9) for quick preset levels (10%-90%)
  - Press 0 for instant 100% brightness
- 🎯 **Centered Layout**: All elements properly centered for a polished look
- 📊 **Visual Feedback**: Progress bar with percentage display
- 🚀 **Lightweight**: Fast and responsive, built with Go
- 🖥️ **Full Screen Mode**: Uses alternate screen buffer (doesn't mess with your terminal history)

## 📋 Requirements

- **Operating System**: Linux
- **Go**: Version 1.21 or higher
- **brightnessctl**: System utility for controlling backlight brightness

### Installing brightnessctl

#### Arch Linux / Manjaro
```bash
sudo pacman -S brightnessctl
```

#### Ubuntu / Debian
```bash
sudo apt install brightnessctl
```

#### Fedora
```bash
sudo dnf install brightnessctl
```

#### From Source
```bash
git clone https://github.com/Hummer12007/brightnessctl
cd brightnessctl
make
sudo make install
```

## 🚀 Installation

### From Source

1. Clone the repository:
```bash
git clone https://github.com/yourusername/brightness-tui.git
cd brightness-tui
```

2. Build the application:
```bash
go build -o brightness-tui
```

3. (Optional) Install to your PATH:
```bash
sudo mv brightness-tui /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/yourusername/brightness-tui@latest
```

## 🎮 Usage

Simply run the application:

```bash
./brightness-tui
```

Or if installed to your PATH:

```bash
brightness-tui
```

### Keybindings

| Key(s) | Action |
|--------|--------|
| `↑` or `k` or `→` or `l` | Increase brightness by 5% |
| `↓` or `j` or `←` or `h` | Decrease brightness by 5% |
| `1` | Set brightness to 10% |
| `2` | Set brightness to 20% |
| `3` | Set brightness to 30% |
| `4` | Set brightness to 40% |
| `5` | Set brightness to 50% |
| `6` | Set brightness to 60% |
| `7` | Set brightness to 70% |
| `8` | Set brightness to 80% |
| `9` | Set brightness to 90% |
| `0` | Set brightness to 100% |
| `q` or `Esc` or `Ctrl+C` | Quit the application |

## 🎨 UI Overview

The interface consists of three main sections:

```
╭─────────────────────────────────────────────────────╮
│                                                     │
│                        󰛨                            │
│                                                     │
│            ██████████████████░░░░░    ← Bar        │
│                     75%               ← Percentage │
│                                                     │
│     ↑/↓: ±5%  •  1-9: 10%-90%  •  0: 100%         │
│              q/esc: quit              ← Controls   │
╰─────────────────────────────────────────────────────╯
```

### Brightness Icon Levels

The brightness icon dynamically changes based on your brightness level:

- **90-100%**: 󰛨 (Maximum brightness)
- **80-89%**: 󱩖
- **70-79%**: 󱩕
- **60-69%**: 󱩔
- **50-59%**: 󱩓
- **40-49%**: 󱩒
- **30-39%**: 󱩑
- **20-29%**: 󱩐
- **10-19%**: 󱩏
- **0-9%**: 󱩎 (Minimum brightness)

## 🛠️ Development

### Project Structure

```
brightness-tui/
├── main.go          # Main application code
├── go.mod           # Go module definition
├── go.sum           # Go module checksums
└── README.md        # This file
```

### Building

```bash
go build -o brightness-tui
```

### Running Tests

```bash
go test ./...
```

### Code Quality

Check for issues:
```bash
go vet ./...
```

Run with staticcheck (if installed):
```bash
staticcheck ./...
```

## 🐛 Troubleshooting

### Permission Denied Error

If you get a permission error when trying to change brightness:

1. Add your user to the `video` group:
```bash
sudo usermod -a -G video $USER
```

2. Create a udev rule for brightnessctl:
```bash
sudo nano /etc/udev/rules.d/90-brightnessctl.rules
```

Add this line:
```
ACTION=="add", SUBSYSTEM=="backlight", RUN+="/bin/chgrp video /sys/class/backlight/%k/brightness"
ACTION=="add", SUBSYSTEM=="backlight", RUN+="/bin/chmod g+w /sys/class/backlight/%k/brightness"
```

3. Reload udev rules:
```bash
sudo udevadm control --reload-rules
sudo udevadm trigger
```

4. Log out and log back in for the group changes to take effect.

### brightnessctl Not Found

If the application can't find `brightnessctl`:

1. Make sure it's installed (see [Requirements](#requirements))
2. Verify it's in your PATH:
```bash
which brightnessctl
```

### Display Issues

If the UI doesn't display correctly:

- Make sure your terminal supports UTF-8
- Try a different terminal emulator (recommended: Alacritty, Kitty, or Wezterm)
- Ensure your terminal size is at least 80x24 characters

## 📦 Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Style definitions for nice terminal layouts

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) by Charm
- Inspired by the need for a better brightness control experience on Linux
- Sun icon concept inspired by various terminal UI projects

## 📧 Contact

For questions, issues, or suggestions, please open an issue on GitHub.

---

Made with ❤️ and Go
