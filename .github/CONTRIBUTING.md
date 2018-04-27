# Contributing to the gosharexserver project

I am glad you are reading this because contributors are always welcomed :).
## Issues

This section deals with the general creation of issues.
### Suggestion of a feature or an enhancement
When suggesting a feature/enhancement, your issue should at least contain the following information:
- short explanation of the feature/enhancement
- (suggested dependencies for realization)
- size of the feature
- importance of the issue
### Report of errors or bugs
When reporting errors or bugs, your issue should at least contain the following information:
- go version (`go version`)
- operating system information (e.g. `ubuntu 16.04 xenial 64bit version`)
- error log
- any additional information of your environment which you think it is important for the issue
Please be aware that normally issues are discussed after they were opened and you may have to provide further information.

## Pull requests

This section deals with the general creation of pull requests.
If you would like to see your code implemented in the main project, there is no way around pull requests. If you open a pull request please structure it like these points describe:
- explanation of the content (or a reference to the issue the pull request deals with)
- list of added dependencies
Before you open a pull request please make sure that you have ran the following commands:
- `go fmt .`
- if you have added new dependencies run `godep ensure` (see [godep](https://github.com/golang/dep/) for more information)

## Code conventions
This application follows the general Golang code convention as well as its commentary conventions. We use simple tabs (`\t`) as the tab intent.
