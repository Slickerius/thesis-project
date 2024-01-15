#ifndef GST_H
#define GST_H

#include <gst/gst.h>
#include <stdint.h>
#include <stdlib.h>

extern void goHandlePipelineBuffer(void *buffer, int bufferLen, int samples, int pipelineId);

void gstreamer_init();
void gstreamer_start_mainloop();

GstElement *gstreamer_receive_create_pipeline(char *pipeline);
void gstreamer_receive_start_pipeline(GstElement *pipeline);
void gstreamer_receive_stop_pipeline(GstElement *pipeline);
void gstreamer_receive_push_buffer(GstElement *pipeline, void *buffer, int len);

GstElement *gstreamer_send_create_pipeline(char *pipeline);
void gstreamer_send_start_pipeline(GstElement *pipeline, int pipelineId);
void gstreamer_send_stop_pipeline(GstElement *pipeline);

void gstreamer_free_pipeline(GstElement *pipeline);

#endif