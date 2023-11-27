# Time Log Synchronization from Link to Jira

This application facilitates the synchronization of time logs from a link to Jira.

## Usage

### Build or Run

You have two options to run the application:

1. Build the binary for your local computer:
   

2. Alternatively, you can directly run the application using:

   ```shell
   go run main.go sync
   ```

### Synchronizing with Specific Date

To synchronize with a specific date, use the following command format:

```shell
go run main.go sync --date <yyyy-mm-dd>
```

For example, to sync logs for August 10, 2023:

```shell
go run main.go sync --date 2023-08-10
```

### Description Format

The application expects a specific format for the description:

```
TRADE-1111:What have you done
```

### Task Selection

After executing the command, in the prompt you can select a particular task to sync or sync everything.

### Duplicate Check

The script checks if you already have a worklog on a task with the same description on the specified day. This should prevent accidental duplicate entries.

## Requirements

- Go (my Golang version go1.20.1 linux/amd64)
