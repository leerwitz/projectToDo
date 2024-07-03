'use strict'

// import { deleteTask } from "./buttonEvents.js";

async function allTaskLoad() {
    try {
        let response = await fetch(`http://localhost:8080/task`);
        
        if (!response.ok) {
            throw new Error(response.statusText);
        }

        let tasks = await response.json();
        let taskList = document.getElementById(`taskList`);
        
        tasks.forEach(task => {
            let taskCurr = document.createElement(`div`);
            let taskCurrNotation;
            
            taskCurr.className = `task`;
            taskCurr.innerHTML = `<h3>Title: ${task.title ? task.title : `No title`} 
            Author: ${task.author}   ID: #${task.id} </h3>`;
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
        });


    } catch(err) {
        console.error(`In function allTaskLoad `, err);
    }
}

async function deleteTask(event) {
    let task = this.closest('div').previousElementSibling;
    let response = await fetch(`http://localhost:8080/task/${task.id}`,{
        method : `DELETE`,
    });

    try {    
        if(!response.ok) {
            throw new Error(`Ошибка` + response.statusText);
        }

        console.log(response.statusText);
        task.remove();
        this.closest('div').remove();
    } catch (err) {
        alert(err.message);
    }


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

async function editTask(event) {
    let task = this.closest('div');
    let taskFields = task.querySelectorAll('p');
    let noEditTask = task.cloneNode(true);
    
    task.querySelector('h2').innerHTML = `
    <div>
    <strong>Title: </strong>
    <input type="text" placeholder="Введите название" pattern="/{2,250}/v" 
    value="${task.querySelector('h2').textContent}" required>
    </div>`;
    taskFields.forEach(field => {
        if(field.id == `Urgent`) {
            field.innerHTML = `
            <div>
                <strong>Urgent: </strong>
                <input class="urgent" name="Urgent" type="checkbox" class="input-text" 
                checked="${field.querySelector('span').textContent == 'Yes' ? true : false}" required>
            </div>`
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

    cancelButton.addEventListener('click', () => {
        task.replaceWith(noEditTask)
        let editButton = noEditTask.querySelector('.edit-button');
        let deleteButton = noEditTask.querySelector(`.close-button`);
        deleteButton.addEventListener('click', deleteTask);
        editButton.addEventListener('click', editTask);
    });

}

document.addEventListener('DOMContentLoaded', function() {
    allTaskLoad();
})