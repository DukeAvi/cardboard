package cards

// Average the area covered by a smaller pixel, so thin details still contribute.
func (r *Renderer) sample(p pixelCanvas, x, y, width, height int) int {
	sourceWidth, sourceHeight := len(p[0]), len(p)
	x0, x1 := x*sourceWidth, (x+1)*sourceWidth
	y0, y1 := y*sourceHeight, (y+1)*sourceHeight
	redSum, greenSum, blueSum, total := 0, 0, 0, 0
	clear := true
	for sy := y0 / height; sy < (y1+height-1)/height; sy++ {
		for sx := x0 / width; sx < (x1+width-1)/width; sx++ {
			area := (min(x1, (sx+1)*width) - max(x0, sx*width)) *
				(min(y1, (sy+1)*height) - max(y0, sy*height))
			index := p[sy][sx]
			if index < 0 {
				index = felt
			} else {
				clear = false
			}
			color := r.colors[index]
			redSum += (color >> 16 & 255) * area
			greenSum += (color >> 8 & 255) * area
			blueSum += (color & 255) * area
			total += area
		}
	}
	if clear {
		return -1
	}
	redValue := (redSum + total/2) / total
	greenValue := (greenSum + total/2) / total
	blueValue := (blueSum + total/2) / total
	color := redValue<<16 | greenValue<<8 | blueValue
	// Values above the eight palette entries hold blended RGB colors.
	return color + 8
}
