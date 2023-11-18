function activateNavItems() {
    const path = window.location.pathname;
    if (path === null || path === undefined) return;

    const e = document.getElementById(path);
    if (e == null) return;
    e.classList.add("nav-item-active");
}
