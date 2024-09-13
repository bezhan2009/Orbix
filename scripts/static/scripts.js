document.querySelector('.menu-toggle').addEventListener('click', function() {
    const menu = document.querySelector('.header-nav-menu');
    menu.classList.toggle('show');
});

document.querySelectorAll('.copy-btn').forEach(button => {
    button.addEventListener('click', () => {
        const codeId = button.getAttribute('data-code');
        const codeElement = document.getElementById(codeId).innerText;
        navigator.clipboard.writeText(codeElement).then(() => {
            button.classList.remove('copy-icon');
            button.classList.add('copied-icon');
            setTimeout(() => {
                button.classList.remove('copied-icon');
                button.classList.add('copy-icon');
            }, 2000);
        });
    });
});

document.addEventListener('DOMContentLoaded', () => {
    const headerNavMenu = document.querySelector('.header');
    let lastScrollTop = 0;

    window.addEventListener('scroll', () => {
        let currentScroll = window.pageYOffset || document.documentElement.scrollTop;

        if (currentScroll > lastScrollTop) {
            headerNavMenu.classList.remove('visible');
            headerNavMenu.classList.add('hidden');
        } else {
            headerNavMenu.classList.remove('hidden');
            headerNavMenu.classList.add('visible');
        }

        lastScrollTop = currentScroll <= 0 ? 0 : currentScroll;
    });})
