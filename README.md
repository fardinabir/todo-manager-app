# ToDo Manager Application

This project is an extension of the [**Fullstack Examination 2024**](https://github.com/zuu-development/fullstack-examination-2024/tree/v0.0.2). The original project provided basic functionalities like adding, editing, viewing, and deleting todo tasks. The following changes were made to address the additional examination tasks:

## 1. Task and Status-Based Filtering
- Introduced filtering functionality based on task name and status.
- Updated the UI to allow users to filter tasks dynamically.
- Modified the `GET /api/v1/todos` endpoint to accept `task` and `status` query parameters.
- Updated the Swagger API documentation to reflect these changes.

## 2. Priority-Based Sorting of Tasks
- Added a **priority** field to the `Todo` model (High, Medium, Low).
- Implemented sorting functionality to prioritize tasks by priority field.
- High-priority tasks are now sorted to the top, followed by medium and low-priority tasks.
- Updated the UI to show tasks list sorted by priority.
- Tests were written for the backend to ensure proper handling of priority.

## 3. Separation of Completed and Incomplete Tasks
- Implemented two separate panels on the UI:
    - **Pending Tasks**: For tasks with a status of `created`.
    - **Completed Tasks**: For tasks with a status of `done`.
- Tasks moved to `done` status are automatically transferred to the "Completed Tasks" panel.
- Tasks marked as incomplete are transferred back to the "Pending Tasks" panel.
- Modified the backend to handle task status updates and reflect them appropriately in the UI.

## 4. Backend API and Test Updates
- Updated backend APIs to handle task filtering and sorting by priority.
- Implemented **payload validations** to ensure that task names, statuses, and priorities are correctly validated.
- Added tests to cover the validation of new sorting, filtering, and task status functionalities.

## 5. Documentation Updates
- Updated the Swagger documentation to reflect changes in API endpoints, including the new task filtering and priority-based sorting functionalities.
- Added detailed descriptions for new parameters and updated the API documentation to be in sync with the changes made.

## Accessing the Application
Once both the backend and frontend servers are running, you can access the application's UI by navigating to the following URL in your web browser: 
[http://localhost:3000](http://localhost:3000). To view the API documentation using Swagger, ensure the backend server is running. Open the following URL in your web browser to access the Swagger UI: [http://localhost:1314/swagger/index.html](http://localhost:1314/swagger/index.html).
For setup and development related instructions please refer to the [original README](./README_OLD.md).
