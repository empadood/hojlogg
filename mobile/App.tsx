import React from 'react';
import { NavigationContainer } from '@react-navigation/native';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { StatusBar } from 'expo-status-bar';
import { HomeScreen } from './src/screens/HomeScreen';
import { NewLogScreen } from './src/screens/NewLogScreen';
import { CameraScreen } from './src/screens/CameraScreen';
import { LogDetailScreen } from './src/screens/LogDetailScreen';

export type RootStackParamList = {
  Home: undefined;
  NewLog: { imageUri?: string; logId?: string } | undefined;
  Camera: { logId?: string } | undefined;
  LogDetail: { id: string };
};

const Stack = createNativeStackNavigator<RootStackParamList>();

export default function App() {
  return (
    <NavigationContainer>
      <StatusBar style="light" />
      <Stack.Navigator
        initialRouteName="Home"
        screenOptions={{
          headerStyle: { backgroundColor: '#1e1e2e' },
          headerTintColor: '#fff',
          headerTitleStyle: { fontWeight: '700' },
          contentStyle: { backgroundColor: '#0f0f1a' },
        }}
      >
        <Stack.Screen name="Home" component={HomeScreen} options={{ title: 'Höjlogg 🏍️' }} />
        <Stack.Screen name="NewLog" component={NewLogScreen} options={{ title: 'New Entry' }} />
        <Stack.Screen
          name="Camera"
          component={CameraScreen}
          options={{ title: 'Scan Dashboard', headerTransparent: true, headerTintColor: '#fff' }}
        />
        <Stack.Screen
          name="LogDetail"
          component={LogDetailScreen}
          options={{ title: 'Log Detail' }}
        />
      </Stack.Navigator>
    </NavigationContainer>
  );
}

