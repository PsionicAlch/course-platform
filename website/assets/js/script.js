document.addEventListener("alpine:init", () => {
  const mobileNav = document.querySelector(".dropdown-navbar");
  if (mobileNav) {
    mobileNav.style = "";
  }

  const courseMenu = document.querySelector(".course-menu");
  if (courseMenu) {
    courseMenu.style = "";
  }

  const adminMenu = document.querySelector(".admin-navbar");
  if (adminMenu) {
    adminMenu.style = "";
  }
});

document.addEventListener("DOMContentLoaded", () => {
  if (typeof hljs !== 'undefined') {
    hljs.highlightAll();
  }
});
