#include "gst.h"

#include <gst/app/gstappsrc.h>

typedef struct SampleHandlerUserData {
  int pipelineId;
} SampleHandlerUserData;

GMainLoop *gstreamer_receive_main_loop = NULL;
void gstreamer_start_mainloop() {
    gstreamer_receive_main_loop = g_main_loop_new(NULL, FALSE);
    g_main_loop_run(gstreamer_receive_main_loop);
}

static gboolean gstreamer_bus_call(GstBus *bus, GstMessage *msg, gpointer data) {
  switch (GST_MESSAGE_TYPE(msg)) {

  case GST_MESSAGE_EOS:
    g_print("End of stream\n");
    break;

  case GST_MESSAGE_ERROR: {
    gchar *debug;
    GError *error;

    gst_message_parse_error(msg, &error, &debug);
    g_free(debug);

    g_printerr("Error: %s\n", error->message);
    g_error_free(error);
    break;
  }
  default:
    break;
  }

  return TRUE;
}

void gstreamer_init() {
    gst_init(NULL, NULL);
}

GstElement *gstreamer_receive_create_pipeline(char *pipeline) {
    GError *error = NULL;
    return gst_parse_launch(pipeline, &error);
}

void gstreamer_receive_start_pipeline(GstElement *pipeline) {
    GstBus *bus = gst_pipeline_get_bus(GST_PIPELINE(pipeline));
    gst_bus_add_watch(bus, gstreamer_bus_call, NULL);
    gst_object_unref(bus);

    gst_element_set_state(pipeline, GST_STATE_PLAYING);
}

void gstreamer_receive_stop_pipeline(GstElement *pipeline) {
    gst_element_set_state(pipeline, GST_STATE_NULL);
}

void gstreamer_receive_push_buffer(GstElement *pipeline, void *buffer, int len) {
    GstElement *src = gst_bin_get_by_name(GST_BIN(pipeline), "src");
    if (src != NULL) {
        gpointer p = g_memdup(buffer, len);
        GstBuffer *buffer = gst_buffer_new_wrapped(p, len);
        gst_app_src_push_buffer(GST_APP_SRC(src), buffer);
        gst_object_unref(src);
    }
}

GstFlowReturn gstreamer_send_new_sample_handler(GstElement *object, gpointer user_data) {
    GstSample *sample = NULL;
    GstBuffer *buffer = NULL;
    gpointer copy = NULL;
    gsize copy_size = 0;
    SampleHandlerUserData *s = (SampleHandlerUserData *)user_data;

    g_signal_emit_by_name (object, "pull-sample", &sample);
    if (sample) {
        buffer = gst_sample_get_buffer(sample);
        if (buffer) {
        gst_buffer_extract_dup(buffer, 0, gst_buffer_get_size(buffer), &copy, &copy_size);
        goHandlePipelineBuffer(copy, copy_size, GST_BUFFER_DURATION(buffer), s->pipelineId);
        }
        gst_sample_unref (sample);
    }

    return GST_FLOW_OK;
}

GstElement *gstreamer_send_create_pipeline(char *pipeline) {
    GError *error = NULL;
    return gst_parse_launch(pipeline, &error);
}

void gstreamer_send_start_pipeline(GstElement *pipeline, int pipelineId) {
    SampleHandlerUserData *s = calloc(1, sizeof(SampleHandlerUserData));
    s->pipelineId = pipelineId;

    GstBus *bus = gst_pipeline_get_bus(GST_PIPELINE(pipeline));
    gst_bus_add_watch(bus, gstreamer_bus_call, NULL);
    gst_object_unref(bus);

    GstElement *appsink = gst_bin_get_by_name(GST_BIN(pipeline), "appsink");
    g_object_set(appsink, "emit-signals", TRUE, NULL);
    g_signal_connect(appsink, "new-sample", G_CALLBACK(gstreamer_send_new_sample_handler), s);
    gst_object_unref(appsink);

    gst_element_set_state(pipeline, GST_STATE_PLAYING);
}

void gstreamer_send_stop_pipeline(GstElement *pipeline) {
    gst_element_set_state(pipeline, GST_STATE_NULL);
}

void gstreamer_free_pipeline(GstElement *pipeline) {
    gst_object_unref(pipeline);
}