const Toast = Swal.mixin({
  toast: true,
  position: "top-end",
  showConfirmButton: false,
  timer: 1000,
  timerProgressBar: true,
  didOpen: (toast) => {
    toast.addEventListener("mouseenter", Swal.stopTimer);
    toast.addEventListener("mouseleave", Swal.resumeTimer);
  },
});

const handleConfirm = (e, options) => {
  e.preventDefault();

  Swal.fire({
    showCloseButton: true,
    showCancelButton: true,
    ...options,
  }).then(({ isConfirmed }) => {
    if (isConfirmed) {
      e.detail.issueRequest();
    }
  });
};

window.Toast = Toast;
window.handleConfirm = handleConfirm;
