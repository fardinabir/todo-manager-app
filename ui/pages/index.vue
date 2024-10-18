<template>
  <div class="todo-main">
    <h1>TODO List | TODO リスト</h1>
    <div v-if="statusMessage" class="status-message">{{ statusMessage }}</div>
    <div class="input-group">
      <input v-model="newTask.task" placeholder="新しいタスクを入力 | Enter a new Task " @keyup.enter="addTodo">
      <select v-model="newTask.priority" @keyup.enter="addTodo">
        <option value=0 disabled selected hidden>Select Priority</option>
        <option value=3>High</option>
        <option value=2>Medium</option>
        <option value=1>Low</option>
      </select>
      <button @click="addTodo"> 追加 </button>
    </div>

    <div class="query-section">
      <input
          v-model="query.task"
          placeholder="Filter by Task Name"
          @input="fetchTodos"
      >
      <select v-model="query.status" @change="fetchTodos">
        <option value="">All Statuses</option>
        <option value="created">Created</option>
        <option value="done">Done</option>
      </select>
    </div>

    <div v-if="hasTodos">
      <div v-if="isPendingView">
        <h2>Pending Tasks</h2>
        <div v-if="pendingTodos.length !== 0" class="task-panel">
          <div v-for="todo in pendingTodos" :key="todo.ID" class="todo-item">
            <input
                v-if="todo.isEditing"
                v-model="todo.Task"
                class="edit-input"
                @blur="editTodo(todo)"
                @keyup.enter="editTodo(todo)"
            >
            <span v-else @click="enableEdit(todo)">{{ todo.Task }}</span>
            <div class="buttons">
              <!-- Priority Button (default) -->
              <template v-if="!todo.isEditingPriority">
                <button class="priority-button"
                        @click="enablePriorityEdit(todo)">
                  {{ getPriorityLabel(todo.Priority) }}
                </button>
              </template>

              <!-- Priority Select (shown when editing) -->
              <template v-else>
                <select v-model="todo.Priority"
                        @blur="updateTodoPriority(todo)"
                        @keyup.enter="updateTodoPriority(todo)">
                  <option value=3>High</option>
                  <option value=2>Medium</option>
                  <option value=1>Low</option>
                </select>
              </template>

              <button @click="updateStatus(todo)">✔️</button>
              <button class="delete-button" @click="deleteTodo(todo.ID)">🗑️</button>
            </div>
          </div>
        </div>
        <div v-else>
          <p>保留中のタスクはありません。| There are no pending tasks.</p>
        </div>
      </div>

      <div v-if="isCompletedView" >
        <h2>Done Tasks</h2>
        <div v-if="completedTodos.length !== 0" class="task-panel">
          <div v-for="todo in completedTodos" :key="todo.ID" class="todo-item">
            <input
                v-if="todo.isEditing"
                v-model="todo.Task"
                class="edit-input"
                @blur="editTodo(todo)"
                @keyup.enter="editTodo(todo)"
            >
            <span v-else class="done-task" @click="enableEdit(todo)">{{ todo.Task }}</span>
            <div class="buttons">
              <button class="done" @click="updateStatus(todo)">✔️</button>
              <button class="delete-button" @click="deleteTodo(todo.ID)">🗑️</button>
            </div>
          </div>
        </div>
        <div v-else>
          <p>完了したタスクはありません。| There are no completed tasks.</p>
        </div>
      </div>
    </div>

    <div v-else>
      <p>タスクがありません。| There are no tasks.</p>
    </div>

  </div>
</template>

<script>
export default {
  data() {
    return {
      newTask: {
        task: '',
        priority: 0,
      },
      todos: [],
      query: {
        task: '',
        status: '',
      },
      statusMessage: '',
    };
  },
  computed: {
    hasTodos() {
      return this.todos.length > 0;
    },
    isPendingView() {
      return !this.query.status || this.query.status === 'created';
    },
    isCompletedView() {
      return !this.query.status || this.query.status === 'done';
    },
    pendingTodos() {
      return this.todos.filter(todo => todo.Status === 'created');
    },
    completedTodos() {
      return this.todos.filter(todo => todo.Status === 'done');
    }
  },
  mounted() {
    this.fetchTodos();
  },
  methods: {
    getPriorityLabel(priority) {
      switch (priority) {
        case 3:
          return 'High';
        case 2:
          return 'Medium';
        case 1:
          return 'Low';
        default:
          return '';
      }
    },
    enablePriorityEdit(todo) {
      todo.isEditingPriority = true;
    },
    async fetchTodos() {
      const params = new URLSearchParams();
      if (this.query.task) params.append('task', this.query.task);
      if (this.query.status) params.append('status', this.query.status);

      const queryString = params.toString() ? `?${params.toString()}` : '';

      this.statusMessage = '';

      try {
        const response = await fetch(`/api/v1/todos${queryString}`, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        });
        if (!response.ok) throw new Error(`Failed to get todo list. statusCode: ${response.status}`);
        response.json().then(data => {
          this.todos = data.data;
          this.statusMessage = `${this.todos.length} task(s) found.`;
        });
      } catch (error) {
        console.error('Error fetching todos:', error);
        this.statusMessage = 'Failed to fetch todos.';
      }
    },
    async updateTodoPriority(todo) {
      todo.isEditingPriority = false;
      todo.Priority = parseInt(todo.Priority, 10);
      try {
        const response = await fetch(`/api/v1/todos/${todo.ID}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            priority: todo.Priority,
          })
        });

        if (!response.ok) throw new Error(`Failed to update priority. statusCode: ${response.status}`);
        this.todos.sort((a, b) => {
          if (a.Priority !== b.Priority) {
            return b.Priority - a.Priority;
          }
          return new Date(b.CreatedAt) - new Date(a.CreatedAt);
        });
        this.statusMessage = '優先度が正常に更新されました| Priority updated successfully';
      } catch (error) {
        console.error('Error updating priority:', error);
        this.statusMessage = '優先順位の更新に失敗しました | Failed to update priority';
      }
    },
    async addTodo() {
      if (this.newTask.task.trim() === '') return;

      try {
        let priority = parseInt(this.newTask.priority, 10) ? parseInt(this.newTask.priority, 10) : 1;

        const response = await fetch('/api/v1/todos', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            task: this.newTask.task,
            priority: priority,
            Status: 'created'
          })
        });

        if (!response.ok) throw new Error(`Failed to create todo. statusCode: ${response.status}`);

        response.json().then(data => {
          this.todos.push(data.data);
          this.todos.sort((a, b) => {
            if (a.Priority !== b.Priority) {
              return b.Priority - a.Priority;
            }
            return new Date(b.CreatedAt) - new Date(a.CreatedAt);
          });
          this.newTask = {
            task: '',
            priority: priority,
          };
          this.statusMessage = 'タスクが追加されました | TODO added';
        });
      } catch (error) {
        console.error('Error creating todo:', error);
        this.statusMessage = 'タスクの作成に失敗しました | Failed To create TODO.';
      }
    },
    enableEdit(todo) {
      todo.isEditing = true;
    },
    async editTodo(todo) {
      todo.isEditing = false;

      try {
        const response = await fetch(`/api/v1/todos/${todo.ID}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            task: todo.Task
          })
        });

        if (!response.ok) throw new Error(`Failed to edit todo. statusCode: ${response.status}`);

        this.statusMessage = 'タスクが編集されました | The task has been edited ' ;
      } catch (error) {
        console.error('Error editing todo:', error);
        this.statusMessage = 'タスクの編集に失敗しました | Failed to edit task';
      }
    },
    async updateStatus(todo) {
      try {
        const response = await fetch(`/api/v1/todos/${todo.ID}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            Status: todo.Status === 'done' ? 'created' : 'done'
          })
        });

        if (!response.ok) throw new Error(`Failed to update todo Status. statusCode: ${response.status}`);

        todo.Status = todo.Status === 'done' ? 'created' : 'done';
        this.statusMessage = 'タスクのステータスが変更されました | The status of the task has changed';
      } catch (error) {
        console.error('Error updating todo Status:', error);
        this.statusMessage = 'ステータスの更新に失敗しました | Failed to update status';
      }
    },
    async deleteTodo(id) {
      try {
        const response = await fetch(`/api/v1/todos/${id}`, {
          method: 'DELETE',
          headers: {
            'Content-Type': 'application/json',
          },
        });

        if (!response.ok) throw new Error(`Failed to update todo Status. statusCode: ${response.status}`);

        this.todos = this.todos.filter(todo => todo.ID !== id);
        this.statusMessage = 'タスクが削除されました | The task has been deleted';
      } catch (error) {
        console.error('Error deleting todo:', error);
        this.statusMessage = 'タスクの削除に失敗しました | Failed to delete task ';
      }
    }
  }
};
</script>

<style scoped>
.query-section {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.query-section input,
.query-section select {
  flex: 1;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.todo-main {
  max-width: 400px;
  margin: 20px auto;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.1);
  background-color: #fff;
}

.task-panel {
  max-width: 800px;
  margin: 5px auto;
  padding: 10px;
  border-radius: 8px;
  box-shadow: 0 0 5px rgba(0, 0, 0, 0.1);
  background-color: #fff;
}

.input-group {
  display: flex;
  margin-bottom: 20px;
}

input, select {
  flex: 1;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  margin-right: 10px;
}

button {
  padding: 10px;
  border: none;
  background-color: #333;
  color: #fff;
  border-radius: 4px;
  cursor: pointer;
}

.todo-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.buttons button {
  background-color: #f0f0f0;
  color: #333;
  margin-left: 5px;
  border-radius: 4px;
  padding: 5px 10px;
  transition: background-color 0.3s, color 0.3s;
}

.buttons button.done {
  background-color: #333;
  color: #fff;
}

.buttons button.done::before {
  color: #fff;
}

.buttons button.delete-button {
  color: white;
}

.status-message {
  margin-bottom: 20px;
  padding: 10px;
  background-color: #f0f0f0;
  border-radius: 4px;
}

.done-task {
  text-decoration: line-through;
}

.edit-input {
  flex: 1;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}
</style>