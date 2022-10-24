module.exports = {
  apps: [{
    name: "twitch-bot",
    script: "just",
    args: "run",
  }],
  deploy: {
    production: {
      user: "root",
      host: "176.126.113.161",
      key: "/home/rprtr258/.ssh/test_vds",
      ref: "origin/master",
      repo: "https://github.com/rprtr258/twitch-bot-api.git",
      path: "/var/www/twitch-bot",
      // TODO: deploy just executable
      // "pre-setup": "go install github.com/rprtr258/rwenv && apt install just", // and install go
      // "post-deploy": "go mod download"
    }
  }
};
