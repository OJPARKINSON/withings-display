from PIL import Image, ImageDraw, ImageFont

def normalize_value(value, min_val, max_val, new_min, new_max):
    """
    Normalize a value from one range to another.
    
    :param value: The value to normalize.
    :param min_val: The current minimum value of the range.
    :param max_val: The current maximum value of the range.
    :param new_min: The new minimum value of the range.
    :param new_max: The new maximum value of the range.
    :return: The normalized value.
    """
    return ((value - min_val) / (max_val - min_val)) * (new_max - new_min) + new_min

def create_dot_graph_image(values, output_path, img_width=400, img_height=200, min_val=None, max_val=None):
    """
    Create an image with a dot graph based on a list of values.

    :param values: List of values to plot.
    :param output_path: Path to save the generated image.
    :param img_width: Width of the image.
    :param img_height: Height of the image.
    :param min_val: Minimum value for normalization (optional).
    :param max_val: Maximum value for normalization (optional).
    """
    # If min_val or max_val are not provided, use the min and max of the values
    if min_val is None:
        min_val = 90
    if max_val is None:
        max_val = 120

    # Create a blank white image
    image = Image.new('RGB', (img_width, img_height), 'white')
    draw = ImageDraw.Draw(image)

    # Define square size
    square_size = 5
    margin = 10

    # Calculate the drawable area within the margins
    drawable_width = img_width - 2 * margin
    drawable_height = img_height - 2 * margin

    # Calculate x coordinates as evenly spaced within the drawable area
    x_coords = [margin + i * (drawable_width // (len(values) - 1)) for i in range(len(values))]

    # Normalize the y values within the drawable area
    y_coords = [normalize_value(value, min_val, max_val, margin + drawable_height, margin) for value in values]

    # Draw lines connecting the points
    for i in range(len(x_coords) - 1):
        draw.line((x_coords[i], y_coords[i], x_coords[i+1], y_coords[i+1]), fill='black')

    # Draw the squares
    for x, y in zip(x_coords, y_coords):
        draw.rectangle((x - square_size, y - square_size, x + square_size, y + square_size), fill='black')

    # Add the max value at the top left and the min value at the bottom right
    # font_size = 10
    font = ImageFont.truetype("arial.ttf", 15)

    max_text = f"Max: {max_val}"
    min_text = f"Min: {min_val}"

    draw.text((margin, margin), max_text, fill="black", font=font)
    draw.text((img_width - margin - font.getsize(min_text)[0], img_height - margin - font_size), min_text, fill="black", font=font)


    # Save the image
    image.save(output_path)
    print(f"Image saved to {output_path}")

# Example values
values = [104, 105, 109, 108, 104, 104, 102, 104, 108, 109]

# Output path for the image
output_path = 'dot_graph.png'

# Create the image
create_dot_graph_image(values, output_path)