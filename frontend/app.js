const apiBaseUrl = 'http://51.250.48.59:8080';
// подставь свой backend URL
let currentUserId = null;

// При загрузке страницы — восстанавливаем user_id из localStorage
window.addEventListener('load', () => {
    const savedUserId = localStorage.getItem('user_id');
    if (savedUserId) {
        currentUserId = savedUserId;
        console.log('Восстановлен user_id:', currentUserId);
        fetchTasks();
    }
});

document.getElementById('logoutButton').addEventListener('click', () => {
    currentUserId = null;
    localStorage.removeItem('user_id'); // удаляем user_id из localStorage
    document.getElementById('taskList').innerHTML = '';
    console.log('Выход выполнен');
});


// Регистрация
function register(username, password) {
    fetch(`${apiBaseUrl}/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ "Name":username, "Password":password })
    })
    .then(response => {
        if (response.ok) {
            alert('Регистрация успешна!');
        } else {
            alert('Ошибка регистрации');
        }
    })
    .catch(error => console.error('Ошибка регистрации:', error));
}

function login(username, password) {
    fetch(`${apiBaseUrl}/login?username=${encodeURIComponent(username)}&password=${encodeURIComponent(password)}`)
        .then(response => response.json())
        .then(data => {
            console.log('Ответ от login:', data);
            if (data.user_id) {
                currentUserId = data.user_id;
                localStorage.setItem('user_id', currentUserId); // сохраняем в localStorage
                console.log('Login success, user_id:', currentUserId);
                fetchTasks();
                alert('Login succeed');
            } else {
                alert('Login failed');
            }
        })
        .catch(error => console.error('Ошибка login:', error));
}


// Получить задачи
function fetchTasks() {
    if (currentUserId === null) {
        console.log('Пользователь не залогинен!');
        return;
    }

    fetch(`${apiBaseUrl}/tasks/${currentUserId}`)
        .then(response => response.json())
        .then(tasks => {
    console.log('Ответ от /tasks:', tasks);
    const taskList = document.getElementById('taskList');
    taskList.innerHTML = '';
    if (Array.isArray(tasks)) {
        tasks.forEach(task => {
            const li = document.createElement('li');
            li.innerHTML = `<strong>${task.header}</strong><br>
            ${task.text}<br>
            <button class="delete-btn" onclick="deleteTask(${task.id})">Удалить</button>`;
            taskList.appendChild(li);
        });
    } else {
        console.error('Ожидался массив задач, но пришло:', tasks);
    }
})
        .catch(error => console.error('Ошибка получения задач:', error));
}

// Добавить задачу
function addTask() {
    if (currentUserId === null) {
        console.log('Пользователь не залогинен!');
        return;
    }

    const taskHeader = document.getElementById('newTaskHeader').value;
    const taskText = document.getElementById('newTaskText').value;

    fetch(`${apiBaseUrl}/tasks`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            header: taskHeader,
            text: taskText,
            userId: currentUserId
        })
    })
    .then(response => {
        if (response.ok) {
            document.getElementById('newTaskHeader').value = '';
            document.getElementById('newTaskText').value = '';
            fetchTasks(); // обновляем список после добавления
        } else {
            alert('Ошибка добавления задачи');
        }
    })
    .catch(error => console.error('Ошибка добавления задачи:', error));
}

function deleteTask(taskId) {
    fetch(`${apiBaseUrl}/tasks/delete/${taskId}`, {
        method: 'DELETE',
    })
    .then(response => {
        if (response.ok) {
            console.log(`Задача ${taskId} удалена`);
            fetchTasks(); // обновляем список
        } else {
            console.error('Ошибка удаления задачи');
        }
    })
    .catch(error => console.error('Ошибка удаления задачи:', error));
}


// Привязка кнопок
document.getElementById('registerButton').addEventListener('click', () => {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    register(username, password);
});

document.getElementById('loginButton').addEventListener('click', () => {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    login(username, password);
});

document.getElementById('refreshButton').addEventListener('click', fetchTasks);

document.getElementById('addTaskButton').addEventListener('click', addTask);
