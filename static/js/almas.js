
      
        document.addEventListener('DOMContentLoaded', function() {
            fetch('/myapp/get_choices')
                .then(response => response.json())
                .then(data => {
                    const carTypes = data.car_types;
                    const fuelTypes = data.fuel_types;
        
                    const carTypeSelect = document.getElementById('carType');
                    const fuelTypeSelect = document.getElementById('fuelType');
        
                    carTypes.forEach(carType => {
                        const option = document.createElement('option');
                        option.value = carType;
                        option.textContent = carType;
                        carTypeSelect.appendChild(option);
                    });
        
                    fuelTypes.forEach(fuelType => {
                        const option = document.createElement('option');
                        option.value = fuelType;
                        option.textContent = fuelType;
                        fuelTypeSelect.appendChild(option);
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
