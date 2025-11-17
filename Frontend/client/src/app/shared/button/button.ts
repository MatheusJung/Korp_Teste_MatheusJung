import { Component ,Input} from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-button',
  imports: [CommonModule],
  templateUrl: './button.html',
  styleUrl: './button.scss',
})
export class Button {
  @Input() title:string = 'Título indisponível';
  @Input() imageUrl:string = 'image/default_photo.svg';
  @Input() linkUrl:string = '';
  @Input() online: boolean = false;
}
