package client.api.mooody.me;

import android.app.Notification;
import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.app.Service;
import android.content.Intent;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.os.Build;
import android.util.Log;
import org.qtproject.qt.android.bindings.QtService;

public class QtAndroidService extends QtService
{
    public native void doWork();
    private static final String CHANNEL_ID = "alarms.notifications.client.api.mooody.me";

    @Override public void onCreate()
    {
        super.onCreate();
        createNotificationChannel();

        Notification notification = new Notification.Builder(this, CHANNEL_ID)
            .setSmallIcon(R.drawable.icon)
            .setContentTitle("MoodyApp Surveillance Service is running")
            .build();
        startForeground(1, notification);
    }

    @Override public void onDestroy()
    {
        super.onDestroy();
    }

    @Override public int onStartCommand(Intent intent, int flags, int startId)
    {
        doWork();
        return START_NOT_STICKY;
    }

    private void createNotificationChannel()
    {
        NotificationChannel serviceChannel = new NotificationChannel(CHANNEL_ID, "Foreground State Pulling Service", NotificationManager.IMPORTANCE_LOW);
        NotificationManager manager = getSystemService(NotificationManager.class);
        manager.createNotificationChannel(serviceChannel);
    }
}
