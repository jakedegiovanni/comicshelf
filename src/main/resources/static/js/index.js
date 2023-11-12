function activateNavItems(id) {
    const e = document.getElementById(id);
    if (e == null) return;
    e.classList.add("nav-item-active");
}
