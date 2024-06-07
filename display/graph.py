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

def draw_rounded_rectangle(draw, xy, radius, fill):
    """
    Draw a rounded rectangle on the image.

    :param draw: The ImageDraw object.
    :param xy: The bounding box (x0, y0, x1, y1).
    :param radius: The radius of the corners.
    :param fill: The fill color.
    """
    x0, y0, x1, y1 = xy
    draw.rectangle([x0 + radius, y0, x1 - radius, y1], fill=fill)
    draw.rectangle([x0, y0 + radius, x1, y1 - radius], fill=fill)
    draw.pieslice([x0, y0, x0 + 2*radius, y0 + 2*radius], 180, 270, fill=fill)
    draw.pieslice([x1 - 2*radius, y0, x1, y0 + 2*radius], 270, 360, fill=fill)
    draw.pieslice([x0, y1 - 2*radius, x0 + 2*radius, y1], 90, 180, fill=fill)
    draw.pieslice([x1 - 2*radius, y1 - 2*radius, x1, y1], 0, 90, fill=fill)


def create_dot_graph_image(values):
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
    img_width = 212
    img_height = 104
    min_val = 90
    max_val = 115

    # Create a blank white image
    image = Image.new('RGB', (img_width, img_height), 'white')
    draw = ImageDraw.Draw(image)

    # Define square size
    square_size = 3
    margin = 10

    # Calculate the drawable area within the margins
    drawable_width = img_width - 2 * margin
    drawable_height = img_height - 2 * margin

    # Calculate x coordinates as evenly spaced within the drawable area
    x_coords = [margin + i * (drawable_width // (len(values) - 1)) for i in range(len(values))]

    # Normalize the y values within the drawable area
    y_coords = [normalize_value(value, min_val, max_val, margin + drawable_height, margin) for value in values]

    # Draw lines connecting the points but not touching the squares
    for i in range(len(x_coords) - 1):
        start_x, start_y = x_coords[i], y_coords[i]
        end_x, end_y = x_coords[i + 1], y_coords[i + 1]
        if start_x != end_x:  # To avoid division by zero if the x coordinates are the same
            slope = (end_y - start_y) / (end_x - start_x)
            intercept = start_y - slope * start_x
            offset = square_size +3  # Offset to prevent lines from touching the squares
            start_x += offset if start_x < end_x else -offset
            end_x -= offset if start_x < end_x else -offset
            start_y = slope * start_x + intercept
            end_y = slope * end_x + intercept
        draw.line((start_x, start_y, end_x, end_y), fill='black')

    # Draw the squares with rounded corners
    for x, y in zip(x_coords, y_coords):
        draw_rounded_rectangle(draw, (x - square_size, y - square_size, x + square_size, y + square_size), radius=square_size//2, fill='black')

    # Add the max value at the top left and the min value at the bottom right
    font_size = 10
    font = ImageFont.load_default()

    max_text = f"{max_val}"
    min_text = f"{min_val}"

    draw.text((margin, margin), max_text, fill="black", font=font)
    draw.text((img_width - margin - font.getlength(min_text), img_height - margin - font_size), min_text, fill="black", font=font)


    # Save the image
    return image


# Example values
weights = [104, 105, 109, 108, 104, 104, 98, 104, 108, 109]

# Output path for the image
output_path = 'dot_graph.png'

# Create the image
image = create_dot_graph_image(weights)

image.save(output_path)
print(f"Image saved to {output_path}")
