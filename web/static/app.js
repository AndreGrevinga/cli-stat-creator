// Keep ONLY drag-and-drop functionality
const fileLabel = document.querySelector(".file-label");
const fileInput = document.getElementById("file-input");
const fileName = document.getElementById("file-name");

// File drag and drop handlers
fileLabel.addEventListener("dragover", (e) => {
  e.preventDefault();
  fileLabel.style.borderColor = "var(--primary-color)";
  fileLabel.style.background = "#eff6ff";
});

fileLabel.addEventListener("dragleave", () => {
  fileLabel.style.borderColor = "var(--border-color)";
  fileLabel.style.background = "var(--bg-secondary)";
});

fileLabel.addEventListener("drop", (e) => {
  e.preventDefault();
  fileLabel.style.borderColor = "var(--border-color)";
  fileLabel.style.background = "var(--bg-secondary)";

  const files = e.dataTransfer.files;
  if (files.length > 0) {
    fileInput.files = files;
    fileName.textContent = files[0].name;
  }
});

// File input change handler
fileInput.addEventListener("change", (e) => {
  if (e.target.files.length > 0) {
    fileName.textContent = e.target.files[0].name;
  }
});

// Optional: htmx event listeners for debugging/logging
document.body.addEventListener("htmx:afterSwap", (e) => {
  console.log("Results loaded successfully");
});
document.body.addEventListener("htmx:responseError", (e) => {
  console.error("Request failed:", e.detail);
});
