
__constant sampler_t sampler = CLK_NORMALIZED_COORDS_FALSE | CLK_ADDRESS_CLAMP_TO_EDGE | CLK_FILTER_NEAREST;

__kernel void invert(__read_only image2d_t src, __write_only image2d_t dest) {
    const int2 pos = {get_global_id(0), get_global_id(1)};
    float4 pixel = read_imagef(src, sampler, pos);

    pixel.r = 1 - pixel.r;
    pixel.g = 1 - pixel.g;
    pixel.b = 1 - pixel.b;

    write_imagef(dest, pos, pixel);
}

