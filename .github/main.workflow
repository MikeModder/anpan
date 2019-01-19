workflow "New workflow" {
  on = "push"
  resolves = ["gofmt"]
}

action "gofmt" {
  uses = "sjkaliski/go-github-actions/fmt@v0.2.0"
  //needs   = "previous-action"
  secrets = ["GITHUB_TOKEN"]
  env = {
    //GO_WORKING_DIR = "./path/to/go/files"
  }
}
