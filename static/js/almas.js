document.addEventListener('DOMContentLoaded', function() {
    fetch('/myapp/get_choices')
        .then(response => response.json())
        .then(data => {
            const carTypes = data.car_types;
            const fuelTypes = data.fuel_types;

            // Assuming you have unique IDs for each car's edit form
            document.querySelectorAll('[id^="editForm-"]').forEach(form => {
                const carTypeSelect = form.querySelector('[name="car_type"]');
                const fuelTypeSelect = form.querySelector('[name="fuel_type"]');
                const currentCarType = carTypeSelect.value;
                const currentFuelType = fuelTypeSelect.value;

                // Clear existing options
                carTypeSelect.innerHTML = '';
                fuelTypeSelect.innerHTML = '';

                // Add options for car types
                carTypes.forEach(carType => {
                    const option = document.createElement('option');
                    option.value = carType;
                    option.textContent = carType;
                    if (carType === currentCarType) {
                        option.selected = true;
                    }
                    carTypeSelect.appendChild(option);
                });

                // Add options for fuel types
                fuelTypes.forEach(fuelType => {
                    const option = document.createElement('option');
                    option.value = fuelType;
                    option.textContent = fuelType;
                    if (fuelType === currentFuelType) {
                        option.selected = true;
                    }
                    fuelTypeSelect.appendChild(option);
                });
            });
        })
        .catch(error => console.error('Error fetching choices:', error));
});












        // Function to add a new file input
        function addImageInput() {
            const container = document.getElementById('imageInputsContainer');
            const newInput = document.createElement('input');
            newInput.type = 'file';
            newInput.name = 'images[]';
            newInput.classList.add('form-control-file', 'image-input');
            newInput.accept = 'image/*';
            container.appendChild(newInput);
        }
    
        // Event listener to handle click on the "Add Another Image" button
        document.getElementById('addImageInput').addEventListener('click', function() {
            addImageInput();
        });
    
        // Event listener to handle change in file inputs
        document.getElementById('imageInputsContainer').addEventListener('change', function(event) {
            if (event.target && event.target.classList.contains('image-input')) {
                const files = event.target.files;
                const previewContainer = document.getElementById('imagePreview');
                for (let i = 0; i < files.length; i++) {
                    const file = files[i];
                    const reader = new FileReader();
                    reader.onload = function() {
                        const image = new Image();
                        image.src = reader.result;
                        image.classList.add('square-image');
                        previewContainer.appendChild(image);
                    }
                    reader.readAsDataURL(file);
                }
            }
        });
    
        // Event listener to handle PDF generation button click
        document.addEventListener("DOMContentLoaded", function() {
            document.getElementById("generateReportButton").addEventListener("click", function() {
                // Send a request to fetch the PDF file
                fetch("/admin/cars/pdf_report")
                    .then(response => response.blob())
                    .then(blob => {
                        // Create a URL for the blob
                        const url = window.URL.createObjectURL(blob);
    
                        // Create a link element
                        const a = document.createElement("a");
                        a.href = url;
                        a.download = "report.pdf"; // Set the file name for download
                        document.body.appendChild(a);
                        
                        // Trigger the click event to start the download
                        a.click();
    
                        // Remove the link element after download
                        window.URL.revokeObjectURL(url);
                        document.body.removeChild(a);
                    })
                    .catch(error => console.error("Error fetching PDF:", error));
            });
        });
