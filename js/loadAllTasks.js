'use strict'

document.addEventListener('DOMContentLoaded', function() {
    allTaskLoad();
})

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
            Author: ${task.author} </h3>`;
            taskCurr.id = task.id;
            taskList.append(taskCurr);

            taskCurrNotation = 
            `<div id="taskModal" class="hidden">
                <h2 id="modalTitle">${task.title ? task.title : `No title`}</h2>
                <p id="modalText">${task.text ? task.text : 'No text'}</p>
                <p id="modalAuthor"><strong>Author:</strong> ${ task.author ? task.author : 'Anonymous'}</p>
                <p id="modalUrgent"><strong>Urgent<strong>: 
                    <span ${task.urgent ? 'style = "color : green">Yes' : 'style= "color : red">No'}</span></p>
            </div>`

            taskCurr.insertAdjacentHTML('afterend', taskCurrNotation);
            taskCurr.addEventListener(`click`, () => {
                let notation = taskCurr.nextSibling;
                    notation.classList.toggle('hidden');
            });
        });


    } catch(err) {
        console.error(`In function allTaskLoad `, err);
    }
}

function showModal(task) {
    document.getElementById('modalTitle').textContent = task?.title;
    document.getElementById('modalText').textContent = task?.text;
    document.getElementById('modalAuthor').textContent = task?.author;
    document.getElementById('modalUrgent').textContent = task?.urgent ? '<span style="color: red">Yes</span>' 
    : '<span style="color: blue">No</span>';

    let modal = document.getElementById('taskModal');
    modal.style.display = "block";

    window.onclick = function(event) {
        if(event.target == modal) {
            modal.style.display = `none`;
        }
    }
}