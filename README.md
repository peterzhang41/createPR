# bitbucket-server-cli-tools
a place for bitbucket server cli tools

## cpr - create pull request

### Release
[cpr v0.1-alpha](https://github.com/peterzhang41/bitbucket-server-go-cli-tools/releases)
    
### Compile
   ~~~
   go install
   ~~~
 cross platform 
   ~~~
   GOOS=windows GOARCH=amd64 go build -o cpr.exe 
   ~~~

### Global Installation

  * Windows: 
    move to or create an folder, get the path and add into PATH variable
    https://www.architectryan.com/2018/03/17/add-to-the-path-on-windows-10/
  * Linux: 
    ~~~
    cp ~/Downloads/cpr /usr/local/bin
    ~~~
 
 ### Example: 
 
  * Set global environment variable pointing to the cpr config file path (Example in Linux ~/.bashrc, windows user, same as the linke above, create a new system variable CPR_CONFIG_FILE_PATH in Environment Variables)
     ~~~
     export CPR_CONFIG_FILE_PATH="/home/cpr_config.yaml"
     ~~~
    then directly execute it , cpr will check $CPR_CONFIG_FILE_PATH
    ~~~
      $ cpr
    ~~~
  * Or, Reading setting from load flag and execute it
      ~~~
        $ cpr --load '/home/cpr_config.yaml'
      ~~~
  * Or, setup|modify flags values
    ~~~
    $ cpr --username fristname.lastname --password '###'\
          --url 'https://bitbucket.example.com' --destBranch 'release/1.2.9'\
          --title 'This is a sample title' --description 'Please check on Line:100'\
          --debug --reviewer firstname.lastname --reviewer firstname.lastname
    ~~~
    
  ### Precedence
  The precedence for flag value sources is as follows (highest to lowest):
  0. Command line flag value from user
  1. Environment variable (if specified)
  2. Configuration file (if specified)
  3. Default defined on the flag
  
  
### Flags
   ~~~
      --load value                    load .yaml config file from the path or from environment variable [$CPR_CONFIG_FILE_PATH]
      --username value                Bitbucket account username
      --password value                Bitbucket account password
      --url value                     bitbucket server url (default: "https://bitbucket.simprocloud.com")
      --destBranch value              PR destination branch
      --title value                   PR title, branch name will be used if the title is not given
      --description value             PR description, could be empty
      --debug                         turn debug on, will turn on all arguments and flags value
      --reviewers firstName.lastName  PR reviewers firstName.lastName
      --help, -h                      show help
      --version, -v                   print the version
  ~~~
