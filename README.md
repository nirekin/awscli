# awscli

**awscli** is a minimalist command line tool helpful to filter a list of AWS EC2 instances and interact with them.

> usage `awscli --help` for the detail of all options.


## Required environment variables

* Set the **AWS_REGION** environment variable to the default region
* Set the **AWS_ACCESS_KEY_ID** and **AWS_SECRET_ACCESS_KEY** environment variables to `<YOUR_ACCESS_KEY_ID>` and `<YOUR_SECRET_ACCESS_KEY>`

## How to filter EC2 instances


Instances will be filtered by the `Tag:Name` content

Usage:
>  `awscli -name <TAG CONTENT>`

It's okay to use willcard `*` or `?` into the tag content.



> `awscli -name testContent*`
> `awscli -name testEnvironment?_testQualifier?`


Note : Without specifying any command it will just returned the list of EC2 intances matching the given `TAG` content.

Example:

```console
> $ ./awscli -name testEnv*
> Working with tagged name:testEnv*
> Located instances matching the tag :3
> testEnvironment4_testQualifier4 i-07789dcc5fa55752b running
> testEnvironment4_testQualifier4 i-0d699bfba90daf8da running
> testEnvironment4_testQualifier4 i-0d7ee2f1c879fc628 running
```

This output shows: 

* The content of the `Tag:Name`
* The instance id
* The instance state
* The public instance IP ( if available )

## Available commands

* **STOP** Command to stop filtered EC2 instances.
* **TERMINATE** Command to terminate filtered EC2 instances.

Examples:

> `awscli -name testContent* -cmd STOP`
> `awscli -name testContent* -cmd TERMINATE`
****

