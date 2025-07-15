#!/usr/bin/env python
from gimpfu import *
import os

def create_new_insta_image(
        padding,
        background_color_r,
        background_color_g,
        background_color_b,
        my_image_path,
        export_path):

    my_image = pdb.gimp_file_load(my_image_path, my_image_path)

    max_dim = max(my_image.width, my_image.height)
    canvas_size = max_dim + padding

    img = gimp.Image(canvas_size, canvas_size, RGB)

    gimp.set_background(background_color_r, background_color_g, background_color_b)

    background = gimp.Layer(img, "Background", canvas_size, canvas_size, RGB_IMAGE, 100, NORMAL_MODE)
    background.fill(BACKGROUND_FILL)
    img.add_layer(background, 0)

    overlay = pdb.gimp_file_load_layer(img, my_image_path)
    img.add_layer(overlay, -1)

    x_offset = (canvas_size - my_image.width) // 2
    y_offset = (canvas_size - my_image.height) // 2
    pdb.gimp_layer_set_offsets(overlay, x_offset, y_offset)

    flattened = pdb.gimp_image_flatten(img)

    ext = os.path.splitext(export_path)[1].lower().strip('.')
    if ext in ["jpg", "jpeg"]:
        pdb.file_jpeg_save(img, flattened, export_path, export_path, 0.9, 0, 0, 0, "", 0, 1, 0, 0)
    
    pdb.gimp_image_delete(my_image)
    pdb.gimp_image_delete(img)

    pdb.gimp_quit(1)

register(
    "python-fu-create_new_insta_image",
    "Create a square Instagram image with padding and chosen background color",
    "Adds padding and background to center an image on a square canvas.",
    "Patryk Ostrowski",
    "Patryk Ostrowski",
    "2025",
    "Create New Insta Image",
    "",
    [
        (PF_INT, "padding", "Padding", 40),
        (PF_INT, "background_color_r", "Red", 0),
        (PF_INT, "background_color_g", "Green", 0),
        (PF_INT, "background_color_b", "Blue", 0),
        (PF_FILE, "my_image_path", "Path to Image", ""),
        (PF_STRING, "export_path", "Export to Path", "")
    ],
    [],
    create_new_insta_image,
    menu="<Image>/Filters/Development/Plug-In Examples/"
)

main()
