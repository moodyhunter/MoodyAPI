package client.api.mooody.me;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.util.Log;

public class QtBootServiceBroadcastReceiver extends BroadcastReceiver
{
    @Override public void onReceive(Context context, Intent intent)
    {
        Intent serviceIntent = new Intent(context, QtAndroidService.class);
        serviceIntent.addFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
        context.startForegroundService(serviceIntent);
    }
}
