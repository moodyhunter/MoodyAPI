package client.api.mooody.me;

import android.app.ActivityManager;
import android.app.ActivityManager.RunningServiceInfo;
import android.app.Notification;
import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.content.Context;
import android.content.Intent;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.graphics.Color;
import android.util.Log;

public class AndroidUtils
{
    public static void notify(Context context, String title, String message)
    {
        NotificationManager m_notificationManager = (NotificationManager) context.getSystemService(Context.NOTIFICATION_SERVICE);

        NotificationChannel notificationChannel;
        notificationChannel = new NotificationChannel("alarm.notifications.client.api.mooody.me", "Camera Alarms", NotificationManager.IMPORTANCE_HIGH);
        m_notificationManager.createNotificationChannel(notificationChannel);

        Bitmap icon = BitmapFactory.decodeResource(context.getResources(), R.drawable.icon);
        Notification.Builder m_builder = new Notification.Builder(context, notificationChannel.getId());
        m_builder.setSmallIcon(R.drawable.icon)
            .setLargeIcon(icon)
            .setContentTitle(title)
            .setContentText(message)
            .setDefaults(Notification.DEFAULT_SOUND);

        m_notificationManager.notify(0, m_builder.build());
    }

    public static void startService(Context context) {
        if (!isMyServiceRunning(context, QtAndroidService.class)) {
            Intent serviceIntent = new Intent(context, QtAndroidService.class);
            serviceIntent.addFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
            context.startForegroundService(serviceIntent);
        }
    }

    private static boolean isMyServiceRunning(Context context, Class<?> serviceClass) {
       ActivityManager manager = (ActivityManager) context.getSystemService(Context.ACTIVITY_SERVICE);
       for (ActivityManager.RunningServiceInfo service : manager.getRunningServices(Integer.MAX_VALUE)) {
          if (serviceClass.getName().equals(service.service.getClassName())) {
             return true;
          }
       }
       return false;
    }
}
