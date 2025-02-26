# YouBlock

YouBlock is your personal productivity sidekick, here to save you from the endless rabbit holes of YouTube, social media, and other digital distractions. Think of it as a bouncer for your brain—keeping the noisy, time-sucking websites out so you can focus on what really matters. Whether you're coding, studying, or just trying to get stuff done.

## Why ?

Let’s face it: the internet is a black hole of distractions. One minute you’re researching something important, and the next, you’re watching a 10-year-old video of a cat playing the piano. YouBlock is here to break the cycle. It’s not just a tool; it’s a lifestyle upgrade. Take back control of your time and attention—because you’ve got better things to do than mindlessly refresh your feed.

## Quick Start

### Installation

1. Clone the repo:
```bash
git clone https://github.com/yourusername/youblock.git
cd youblock
```

2. Run the installer:
```bash
sudo ./install-youblock
```

3. Sit back and let YouBlock do its thing. It’s like hiring a bodyguard for your productivity.

### Uninstallation

If you ever want to break free (we won’t judge), just run:

```bash
sudo ./uninstall-youblock
```

## Usage

### How It Works

YouBlock quietly runs in the background, keeping your /etc/hosts file in check. If anyone (including you) tries to unblock a distracting site, YouBlock steps in and says, “Not so fast!” It’s like having a productivity ninja on your team.

### Managing the Service

```bash
systemctl status youblock
```

```bash
journalctl -u youblock.service
Restart the Service:
```

```bash
sudo systemctl restart youblock
```

## Contributing

YouBlock is a community-driven project, and we’d love your help! Whether you’re a Go guru, a documentation wizard, or just someone with great ideas, there’s a place for you here.

## License

YouBlock is open-source and free under the MIT License. Use it, tweak it, share it—just don’t blame us if you suddenly become way more productive. 😉
