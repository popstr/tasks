# Tasks

## API Endpoints

`GET /tasks` - list all tasks

Query params:

- `category` - filter tasks by category name
- `sortBy` - sort by field (`name`, `category`, `dueDate`, `status`)
- `sortDir` - sorting direction, asc or desc

Example: /tasks?category=Category%203&sortBy=name&sortDir=desc

`GET /tasks/[taskId]` - get a single task by id

Example: `/tasks/1`

```
{
    "id": 1,
    "name": "Do the dishes",
    "category": "Household",
    "dueDate": "2023-08-14",
    "status": "todo"
}
```

`POST /tasks` - Create new task.

Example:

```
{
    "name": "The new task",
    "category": "New",
    "dueDate": "2023-06-30",
    "status": "todo"
}
```

`PUT /tasks/[taskId]` - Update task

You only need to include the fields you want to update.

Example: `/tasks/1`

```
{
    "category": "Completed",
    "status": "done"
}
```

`DELETE /tasks/[taskId]` - Delete single task

## Additional information

Due date is given in the format YYYY-MM-DD, e.g. 2023-08-24.

Status can be `todo`, `doing` or `done`.
