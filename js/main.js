'use strict'

const serverURL = `http://localhost:8080/`
const createTaskLocation = `http://localhost/html/createTask.html`

async function allTaskLoad() {
    try {
        let response = await fetch(serverURL + `task`);
        
        if (!response.ok) {
            throw new Error(response.statusText);
        }

        let tasks = await response.json();
        let taskList = document.getElementById(`taskList`);
        
        tasks.forEach(task => addTaskToTaskList(task));

    } catch(err) {
        console.error(`In function allTaskLoad `, err);
    }
}

async function deleteTask(event) {
    let task = event.currentTarget.closest('div').previousElementSibling;
    let response = await fetch(serverURL + `task/${task.id}`,{
        method : `DELETE`,
    });

    try {    
        if(!response.ok) {
            throw new Error(`Ошибка` + response.statusText);
        }

        task.remove();
        this.closest('div').remove();
    } catch (err) {
        alert(err.message);
    }


}

function addTaskToTaskList(task) {
    let taskCurr = document.createElement(`div`);
    let taskCurrNotation;
    
    taskCurr.className = `task`;
    taskCurr.innerHTML = `
    <h3>Title: ${task.title ? task.title : `No title`} 
    Author: ${task.author}
    Urgent: <span ${task.urgent ? 'style = "color : green">Yes' : 'style= "color : red">No'}</span>    
    ID: #${task.id}</h3>`;
    taskCurr.id = task.id;
    taskList.append(taskCurr);

    taskCurrNotation = getStrTaskFromObj(task);

    taskCurr.insertAdjacentHTML('afterend', taskCurrNotation);
    taskCurr.addEventListener(`click`, () => {
        let notation = taskCurr.nextSibling;
            notation.classList.toggle('hidden');
    });

    let editButton = taskCurr.nextElementSibling.querySelector('.edit-button');
    let deleteButton = taskCurr.nextElementSibling.querySelector(`.close-button`);
    deleteButton.addEventListener('click', deleteTask);
    editButton.addEventListener('click', editTask);
}

function getTaskFromDiv(div) {
    let task = {};
    
    for(let child of div.children) {
        task[child.name] = child.textContent;
    }

    return task;
}

function getStrTaskFromObj(task) {
    return `<div id="taskModal" class="hidden" style="position : relative">
                <h2 id="title">${task.title ? task.title : `No title`}</h2>
                <p id="Text">${task.text ? task.text : 'No text'}</p>
                <p id="Author"><strong>Author:</strong> ${ task.author ? task.author : 'Anonymous'}</p>
                <p id="Urgent"><strong>Urgent<strong>: 
                <span ${task.urgent ? 'style = "color : green">Yes' : 'style= "color : red">No'}</span>
                <button class="edit-button">Редактировать</button> <button class="close-button right-bottom">Удалить</button></p>
        
            </div>`
}

function cancelEditTask(task, noEditTask) {
    task.innerHTML = noEditTask
}

function replaceTask(task, newTask) {
    return function () {
        task.replaceWith(newTask)
        let editButton = newTask.querySelector('.edit-button');
        let deleteButton = newTask.querySelector(`.close-button`);
        deleteButton.addEventListener('click', deleteTask);
        editButton.addEventListener('click', editTask);
    }
}

function editButtonEvent(task) {

    return async function(event) {
        let inputs = task.querySelectorAll('input');
        let newTask = {};
        let headerTask = task.previousElementSibling;
        try {
            inputs.forEach(input => {
                if(input.type == 'checkbox'){
                    newTask[input.name[0].toLowerCase() + input.name.slice(1)] = input.checked;
                } else {
                    newTask[input.name[0].toLowerCase() + input.name.slice(1)] = input.value;
                }
            });

            let response = await fetch(serverURL + `task/${task.previousElementSibling.id}`, {
                method : 'PATCH',
                headers : {
                    'Content-Type' : 'applications/json'
                },
                body : JSON.stringify(newTask),
            });

            if(!response.ok) {
                throw new Error("Ошибка: " + response.statusText);
            }

            newTask.id = task.previousElementSibling.id; 
            task.insertAdjacentHTML('afterend', getStrTaskFromObj(newTask));
            task.nextElementSibling.classList.remove('hidden');
            task.nextElementSibling.querySelector(`.edit-button`).addEventListener(`click`, editTask);
            task.nextElementSibling.querySelector(`.close-button`).addEventListener(`click`, deleteTask);
            task.replaceWith(task.nextElementSibling);
            headerTask.innerHTML = `
                <h3>Title: ${newTask.title ? newTask.title : `No title`} 
                Author: ${newTask.author} 
                Urgent: <span ${task.urgent ? 'style = "color : green">Yes' : 'style= "color : red">No'}</span>    
                ID: #${newTask.id}</h3>`;

        } catch(err) {
            alert(err.message);
        } 
    }
}

async function editTask(event) {
    let task = this.closest('div');
    let taskFields = task.querySelectorAll('p');
    let noEditTask = task.cloneNode(true);
    let isChecked = task.querySelector('span').textContent == 'Yes';

    task.querySelector('h2').innerHTML = `
    <div>
    <strong>Title: </strong>
    <input name="title" type="text" placeholder="Введите название" pattern="/{2,250}/v" 
    value="${task.querySelector('h2').textContent}" required>
    </div>`;
    
    taskFields.forEach(field => {
        if(field.id == `Urgent`) {
        
            field.innerHTML = `
            <div>
                <strong>Urgent: </strong>
                <input class="urgent" name="Urgent" type="checkbox" class="input-text" required>
            </div>`
            field.querySelector(`input`).checked = isChecked;
            
        } else {
        
            field.innerHTML = `<div>
            <strong>${field.id}: </strong>
            <input name=${field.id} type="text" class="input-text" 
            value="${field.lastChild.textContent}" required>
            </div>`
        }
    });

    let buttons = `<button class="edit-button">Редактировать</button> <button class="close-button right-bottom">Отмена</button>`
    task.lastChild.innerHTML += buttons;

    let cancelButton = task.querySelector(`.close-button`);
    let editButton = task.querySelector('.edit-button');

    cancelButton.addEventListener('click', replaceTask(task, noEditTask));

    editButton.addEventListener('click', editButtonEvent(task));

}

async function searchTaskEvent(event) {
    let input = this.previousElementSibling.value;
    let taskList = document.getElementById(`taskList`);
    try {

        if(input[0] == '#') {
            let response =  await fetch(serverURL + `task/${input.slice(1)}`);
            
            if(!response.ok) {
                throw new Error("Ошибка: " + response.statusText)
            }

            let result = await response.json();
            taskList.innerHTML = '';
            addTaskToTaskList(result);

        } else {
            let url = new URL(serverURL + 'task');
            url.searchParams.append(`title`, input);

            let response = await fetch(url);

            if(!response.ok) {
                throw new Error(`Ошибка: ` + response.statusText);
            }

            let result =  await response.json();
            taskList.innerHTML = '';

            result.forEach(task => addTaskToTaskList(task));
        }

    } catch(err) {
        console.error(err.message);
    }
}

document.addEventListener('DOMContentLoaded', function() {
    let searchButton = document.getElementById('searchButton');
    let createButton = document.getElementById('createButton');

    allTaskLoad();
    createButton.addEventListener('click', function() {
      window.location.href = createTaskLocation;
    });
    
    searchButton.addEventListener('click', searchTaskEvent);

})