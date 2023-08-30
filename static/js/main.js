function copyAnswerToClipboard() {
    const btn = document.querySelector('#copy-button')
    const answer = document.querySelector('#answer');
    const color = getComputedStyle(document.body).getPropertyValue('--green');

    navigator.clipboard.writeText(answer.href)
        .then(() => {
            btn.innerText = 'Done';
            btn.style.color = 'black';
            btn.style.fontWeight = 'bold';
            btn.style.backgroundColor = color;
            btn.style.borderColor = color;
        })
        .catch(err => console.error(err))
}

function main() {
    const btn = document.querySelector('#copy-button');
    btn.addEventListener('click', copyAnswerToClipboard)
}

document.addEventListener('DOMContentLoaded', main)