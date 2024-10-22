// Function to generate gradient colors
function generateGradientColors(startColor, endColor, steps) {
  const start = hexToRgb(startColor);
  const end = hexToRgb(endColor);
  const stepFactor = 1 / (steps - 1);
  const gradient = [];

  for (let i = 0; i < steps; i++) {
    gradient.push(interpolateColor(start, end, stepFactor * i));
  }

  return gradient.map(rgbToHex);
}

// Convert hex color to RGB
function hexToRgb(hex) {
  let bigint = parseInt(hex.replace("#", ""), 16);
  let r = (bigint >> 16) & 255;
  let g = (bigint >> 8) & 255;
  let b = bigint & 255;
  return [r, g, b];
}

// Convert RGB color to hex
function rgbToHex(rgb) {
  return "#" + rgb.map(x => x.toString(16).padStart(2, "0")).join("");
}

// Interpolate between two colors
function interpolateColor(color1, color2, factor) {
  if (arguments.length < 3) factor = 0.5;
  let result = color1.slice();
  for (let i = 0; i < 3; i++) {
    result[i] = Math.round(result[i] + factor * (color2[i] - color1[i]));
  }
  return result;
}

// Apply gradient to each li in the given ul
function applyGradientToLiItems(ul, startColor, endColor) {
  const lis = ul.querySelectorAll("li");
  const colors = generateGradientColors(startColor, endColor, lis.length);

  lis.forEach((li, index) => {
    li.style.paddingLeft = "0.5rem";
    li.style.borderLeftWidth = "5px";
    li.style.borderLeftStyle = "solid";
    li.style.borderLeftColor = colors[index];
    li.style.marginTop = "1rem";
    if (index === 0) li.style.marginTop = "0"; // No margin for the first li
  });
}

// Main function to apply gradient to all ul elements with class 'gradient-list'
function applyGradientToAllLists() {
  const uls = document.querySelectorAll("ul.gradient-list");

  uls.forEach(ul => {
    applyGradientToLiItems(ul, "#00E1FF", "#FF1E00");
  });
}

// Export the main function for external use
export { applyGradientToAllLists };
