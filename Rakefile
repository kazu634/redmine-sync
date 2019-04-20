desc 'Install the initial `go` binaries'
task :setup do
  sh 'go get github.com/golang/dep/cmd/dep'
  sh 'go get github.com/golang/lint/golint'
  sh 'go get golang.org/x/tools/cmd/goimports'
end

desc 'Install dependencies'
task deps: [ :setup ] do
  sh 'env DEPNOLOCK=1 dep init'
  sh 'env DEPNOLOCK=1 dep ensure'
end

desc 'Format source codes'
task fmt: [ :setup ] do
  sh 'find `pwd` -type d -name lib -prune -o -type f -name "*.go" | xargs -t -n 1 goimports -w '
end

desc 'Lint'
task lint: [ :setup ] do
  `glide novendor -x`.split.each do |target|
    sh "golint -set_exit_status #{target} || exit $?"
  end
end

desc 'Build binary'
task :build do
  sh 'git rev-parse --is-inside-work-tree' do |ok, status|
    if ok
      version = `git describe --tags --abbrev=0`.chomp
      revision = `git rev-parse --short HEAD`.chomp
    else
      version = '0.0'
      revision = 'xxxxxxxx'
    end

    ldflags = "-X main.version=#{version} -X main.revision=#{revision}"

    sh "go build -ldflags \"#{ldflags}\" main.go"
  end
end
