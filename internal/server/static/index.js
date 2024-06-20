function activateNavItems() {
    const e = document.getElementById(window.location.pathname);
    if (e == null) return;
    e.classList.add("nav-item-active");
}
