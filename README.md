# Bitbucket Server Cli Tools
a collection of bitbucket server cli tools

## CPR - a tool for creating pull request in Bitbucket server

### Release
[cpr v0.1-alpha](https://github.com/peterzhang41/bitbucket-server-go-cli-tools/releases)
    
#### you could generate binary file by yourself. Requires Golang Environment
   ~~~
   go install
   ~~~
 cross platform 
   ~~~
   GOOS=windows GOARCH=amd64 go build -o cpr.exe 
   ~~~

### Global Installation
  * Download sample-cpr-config.yaml from release, and modify it according to your own config
  * Download cpr_XXX_XXX binary file and rename to cpr.exe (windows) or cpr （linux）
  * Windows:  
    move binary file to a folder which the path needs to be added or has been in PATH variable  
    https://www.architectryan.com/2018/03/17/add-to-the-path-on-windows-10/  
    (optional)  
    you could create a new variable CPR_CONFIG_FILE_PATH and value is the config file path
  * Linux:  
    copy binary to bin folder
    ~~~
    cp ~/Downloads/cpr /usr/local/bin
    ~~~
    add one more line in  ~/.bashrc (modify other init file if you are not using bash)
     ~~~
     export CPR_CONFIG_FILE_PATH="/home/cpr_config.yaml"
     ~~~
    save and execute below, or relaunch terminal
     ~~~
     source ~/.bashrc
     ~~~
  * Check success
     ~~~
     cpr -h
     ~~~
 
 ### Example: 
 
  * you can directly execute it , cpr will check $CPR_CONFIG_FILE_PATH
    ~~~
      cpr
    ~~~
  * Or, execute it by loading settings from the path defined on the flag
      ~~~
      cpr --load '/home/cpr_config.yaml'
      ~~~
  * Or, setup|modify flags values in CLI
    ~~~
    cpr --username fristname.lastname --password '###'\
          --url 'https://bitbucket.example.com' --destBranch 'release/1.2.9'\
          --title 'This is a sample title' --description 'Please check on Line:100'\
          --debug --reviewer firstname.lastname --reviewer firstname.lastname
    ~~~
  * An example below created a PR to 'release/1.2.9' (use single quote if has special character), and used all other default settings in yaml file if it has been setup correctly.  
    ~~~
    cpr --destBranch 'release/1.2.9' --description 'Please check on Line:100'
    ~~~
    
  ### Precedence
  The precedence for flag value sources is as follows (highest to lowest):  
  4. Command line flag value from user  
  3. Environment variable (if specified)  
  2. Configuration file (if specified)  
  1. Default defined on the flag  
  
  
### Flags
    (--description and --debug is not configurable in yaml file, CLI only )  
   ~~~
      --load path                     load .yaml config file from the path or from environment variable [$CPR_CONFIG_FILE_PATH]
      --username firstName.lastName   Bitbucket account username
      --password '#@#$%***'           Bitbucket account password
      --url value                     bitbucket server url (default: "https://bitbucket.simprocloud.com")
      --destBranch value              PR destination branch
      --title value                   PR title, branch name will be used if the title is not given
      --description value             PR description, could be empty
      --debug                         turn debug on will print out all arguments
      --reviewers firstName.lastName  PR reviewers firstName.lastName, it could be multiple
      --help, -h                      show help
      --version, -v                   print the version
  ~~~
