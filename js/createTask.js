'use strict'




async function getTaskFromForm(event) {
    let table = this.closest('table');
    let form = {};
    try{
        for (let field of table.querySelectorAll('input[name]')) {
            if (field.type == 'text') {
                form[`${field.name}`] = field.value;
            } else if (field.type == 'checkbox') {
                form[`${field.name}`] = field.checked;
            }
        }

        let formJSON = JSON.stringify(form);
        let response = await fetch('http://localhost:8080/task', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: formJSON,
        });

        if (!response.ok) {
            throw new Error('Ошибка: ' + response.statusText);
        }

        let result = await response.json();


        table.remove();
        let headerSuccessCreate = document.createElement('h1');
        headerSuccessCreate.textContent = "Новая задача успешно создана!";
        headerSuccessCreate.className = 'success-header center';
        let newTask = 
        `<div id="taskModal" class=" task center">
                <h2 id="modalTitle">${result.title ? result.title : `No title`}</h2>
                <p id="modalText">${result.text ? result.text : 'No text'}</p>
                <p id="modalAuthor"><strong>Author:</strong> ${ result.author ? result.author : 'Anonymous'}</p>
                <p id="modalUrgent"><strong>Urgent<strong>: 
                    <span ${result.urgent ? 'style = "color : green">Yes' : 'style= "color : red">No'}</span></p>
                    <buton id="escapeButton" class="edit-button" style="right : 10px">На главную</buton>
            </div>`;

        document.body.append(headerSuccessCreate);
        headerSuccessCreate.insertAdjacentHTML('afterend', newTask);
        document.getElementById('escapeButton').addEventListener('click', () => window.history.back());
        setTimeout(() => window.history.back(), 5000);

    } catch(err) {
        alert(err.message);
    }
}

// export async function deleteTask(event) {
//     let task = this.closest('div').previousElementSibling;
//     let response = await fetch(`http://localhost:8080/task/${task.id}`,{
//         method : `DELETE`,
//     });

//     try {    
//         if(!response.ok) {
//             throw new Error(`Ошибка` + response.statusText);
//         }

//         console.log(response.statusText);
//         task.remove();
//         this.closest('div').remove();
//     } catch (err) {
//         alert(err.message);
//     }


// }


let createButton = document.getElementById('createButton');
let closeButtons = document.querySelectorAll('.close-button');

closeButtons.forEach(button => {
    button.addEventListener('click', function() {
        window.history.back();
    });        
});

createButton.addEventListener('click',getTaskFromForm);
